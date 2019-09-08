/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trace

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"

	"k8s.io/klog"
)

const (
	traceIDHeader = "x-kubernetes-trace-id"
	spanIDHeader  = "x-kubernetes-span-id"
)

var (
	// DebugAssert causes tracing functionality to panic in various recoverable
	// scenarios (e.g. when an unimplemented method is called, finish called
	// twice, span methods called from multiple goroutines). Setting this value
	// to true is useful in tests.
	DebugAssert = false

	noopTracer opentracing.Tracer = opentracing.NoopTracer{}

	now = time.Now
)

type kTracer struct {
	sampler Sampler
	pool    pool
}

var _ = opentracing.Tracer(&kTracer{})

// StartSpan creates new Span with the given name and options.
//
// See https://godoc.org/github.com/opentracing/opentracing-go#Tracer
func (t *kTracer) StartSpan(name string, opts ...opentracing.StartSpanOption) opentracing.Span {
	os := opentracing.StartSpanOptions{}
	for _, o := range opts {
		o.Apply(&os)
	}

	var traceID, spanID uint64
	if len(os.References) == 0 {
		// This is a root span so check whether we should sample it.
		traceID = newID()
		if !t.sampler.Sample(traceID) {
			return noopTracer.StartSpan("")
		}
		// Set the root span ID to the trace ID for simplicity.
		spanID = traceID
	} else {
		spanID = newID()
	}

	if len(os.References) > 1 {
		maybePanic("wat? when would this happen")
	}

	var parentSpanID uint64
	for _, ref := range os.References {
		switch ref.Type {
		case opentracing.ChildOfRef, opentracing.FollowsFromRef:
			parentSpan := ref.ReferencedContext.(kSpanContext)
			traceID = parentSpan.traceID
			parentSpanID = parentSpan.spanID
		}
	}

	start := os.StartTime
	if os.StartTime.IsZero() {
		start = now()
	}

	s := t.pool.get()
	s.spanID = spanID
	s.parentSpanID = parentSpanID
	s.traceID = traceID
	s.tracer = t
	s.start = start
	s.finish = time.Time{}
	s.drain = t.pool.put
	return s
}

// Inject injects a the SpanContext into the carrier.
//
// See https://godoc.org/github.com/opentracing/opentracing-go#Tracer
func (*kTracer) Inject(spanCtx opentracing.SpanContext, format interface{}, carrier interface{}) error {
	ksc := spanCtx.(kSpanContext)
	switch format {
	case opentracing.HTTPHeaders:
		hc := carrier.(opentracing.HTTPHeadersCarrier)
		hc.Set(traceIDHeader, toStringID(ksc.traceID))
		hc.Set(spanIDHeader, toStringID(ksc.spanID))
	default:
		// TextMap and Binary carriers are not supported.
		maybePanic("unknown carrier")
	}
	return nil
}

// Extract extracts a the SpanContext out of the carrier the carrier.
//
// See https://godoc.org/github.com/opentracing/opentracing-go#Tracer
func (*kTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	switch format {
	case opentracing.HTTPHeaders:
		hc := carrier.(opentracing.HTTPHeadersCarrier)
		var spanCtx kSpanContext
		if err := hc.ForeachKey(func(k, v string) error {
			var err error
			if strings.EqualFold(k, traceIDHeader) {
				spanCtx.traceID, err = fromStringID(v)
				return err
			}
			if strings.EqualFold(k, spanIDHeader) {
				spanCtx.spanID, err = fromStringID(v)
				return nil
			}
			return nil
		}); err != nil {
			return nil, err
		}
		if spanCtx.traceID == 0 {
			return nil, opentracing.ErrSpanContextNotFound
		}
		if spanCtx.spanID == 0 {
			// TODO(mikedanese): document why this is a good think to do.
			spanCtx.spanID = spanCtx.traceID
		}
		return &spanCtx, nil
	default:
		// TextMap and Binary carriers are not supported.
		panic("unsupported carrier")
	}
}

type kSpan struct {
	spanID, parentSpanID, traceID uint64
	tracer                        opentracing.Tracer
	start, finish                 time.Time
	logs                          []opentracing.LogRecord
	drain                         func(span *kSpan)
}

var _ = opentracing.Span(&kSpan{})

// Sets the end timestamp and finalizes Span state.
//
// With the exception of calls to Context() (which are always allowed),
// Finish() must be the last call made to any span instance, and to do
// otherwise leads to undefined behavior.
func (s *kSpan) Finish() {
	s.FinishWithOptions(opentracing.FinishOptions{})
}

// FinishWithOptions is like Finish() but with explicit control over
// timestamps and log data.
func (s *kSpan) FinishWithOptions(opts opentracing.FinishOptions) {
	s.finish = opts.FinishTime
	if s.finish.IsZero() {
		s.finish = now()
	}
	s.drain(s)
}

// Context() returns the SpanContext.
func (s *kSpan) Context() opentracing.SpanContext {
	return kSpanContext{
		traceID: s.traceID,
		spanID:  s.spanID,
	}
}

// LogFields
func (s *kSpan) LogFields(fields ...opentracinglog.Field) {
	s.logs = append(s.logs, opentracing.LogRecord{
		Timestamp: now(),
		Fields:    fields,
	})
}

// LogKV
func (s *kSpan) LogKV(alternatingKeyValues ...interface{}) {
	fields, err := opentracinglog.InterleavedKVToFields(alternatingKeyValues)
	if err != nil {
		s.LogFields(opentracinglog.Error(fmt.Errorf("failed to LogKV: %v", err)))
		return
	}
	s.LogFields(fields...)
}

// Tracer provides access to the Tracer that created this Span.
func (s *kSpan) Tracer() opentracing.Tracer {
	return s.tracer
}

// SetOperationName is unimplemented. Operation names of Span are immutable but
// this method is required to implement the opentracing.Span interface.
func (s *kSpan) SetOperationName(operationName string) opentracing.Span {
	maybePanic("unimplemented: see godoc")
	return s
}

// SetTag is unimplemented.
func (s *kSpan) SetTag(key string, value interface{}) opentracing.Span {
	maybePanic("unimplemented: see godoc")
	return s
}

// SetBaggageItem is unimplemented. It's implementation is costly so it should
// remain unimplemented until we need it
func (s *kSpan) SetBaggageItem(k, v string) opentracing.Span {
	maybePanic("unimplemented: see godoc for Span.SetBaggageItem")
	return s
}

// BaggageItem is unimplement. See Span.SetBaggageItem
func (s *kSpan) BaggageItem(string) string {
	maybePanic("unimplemented: see godoc for Span.SetBaggageItem")
	return ""
}

// Log is unimplemented. Use LogFields or LogKV.
func (*kSpan) Log(opentracing.LogData) {
	maybePanic("unimplemented: use LogFields or LogKV")
}

// LogEvent is unimplemented. Use LogFields or LogKV.
func (*kSpan) LogEvent(string) {
	maybePanic("unimplemented: use LogFields or LogKV")
}

// LogEventWithPayload is unimplemented. Use LogFields or LogKV.
func (*kSpan) LogEventWithPayload(string, interface{}) {
	maybePanic("unimplemented: use LogFields or LogKV")
}

type kSpanContext struct {
	traceID uint64
	spanID  uint64
}

func (kSpanContext) ForeachBaggageItem(func(k, v string) bool) {
	maybePanic("unimplemented: see godoc for Span.SetBaggageItem")
}

func maybePanic(msg string) {
	if DebugAssert {
		panic(msg)
	}
	klog.Error(msg)
}

// NewSpan takes an existing context.Context, creates a new trace.Span and
// returns a new context.Context to be used to propagate the Span (wrapping the
// passed in Context, if it was passed in) and the newly created Span.
//
// NOTE: Be sure to call "span.Finish()" to have the Span logged.
func NewSpan(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span) {
	tracer := noopTracer
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		tracer = parentSpan.Tracer()
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span := tracer.StartSpan(name, opts...)
	return opentracing.ContextWithSpan(ctx, span), span
}

func newID() uint64 {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	return binary.BigEndian.Uint64(buf)
}

const hexStringUint64Length = 18

var hexStringPrefix = []byte{'0', 'x'}

func toStringID(id uint64) string {
	buf := make([]byte, hexStringUint64Length)
	copy(buf[0:2], hexStringPrefix)

	ibuf := make([]byte, 8)
	binary.BigEndian.PutUint64(ibuf, id)

	hex.Encode(buf[2:], ibuf)
	return string(buf)
}

func fromStringID(id string) (uint64, error) {
	b := []byte(id)
	if len(b) != hexStringUint64Length {
		return 0, fmt.Errorf("unexpected hex uint64 length %d: %s", len(b), id)
	}
	if bytes.Compare(hexStringPrefix, b[0:2]) != 0 {
		return 0, fmt.Errorf("unexpected hex uint64 fromat: %s", id)
	}

	ibuf := make([]byte, 8)
	_, err := hex.Decode(ibuf, b[2:])
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(ibuf), nil
}

type Sampler interface {
	Sample(traceID uint64) bool
}

type maskSampler uint64

func (m maskSampler) Sample(traceID uint64) bool {
	return uint64(m)&traceID != uint64(m)
}

var Always Sampler = maskSampler(0xFFFFFFFFFFFFFFFF)
