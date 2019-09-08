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
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
)

func TestMain(m *testing.M) {
	flag.BoolVar(&DebugAssert, "tracing-assert", true, "Enable tracing debug assertions")
	flag.Parse()
	os.Exit(m.Run())
}

func makeTree(n *traceNode, width, depth int) {
	if depth == 0 {
		return
	}
	for i := 0; i < width; i++ {
		child := &traceNode{}
		makeTree(child, width, depth-1)
		n.children = append(n.children, child)
	}
}

func BenchmarkDeepTreeNoPool(b *testing.B) {
	benchmarkDeepTree(b, func() *kTracer {
		return &kTracer{
			sampler: Always,
			pool:    &noopSpanPool{},
		}
	})
}

func BenchmarkDeepTreePooled(b *testing.B) {
	benchmarkDeepTree(b, func() *kTracer {
		pool := newSpanPool(10000)
		return &kTracer{
			sampler: Always,
			pool:    pool,
		}
	})
}

func benchmarkDeepTree(b *testing.B, newTracer func() *kTracer) {
	cases := []struct {
		width, depth int
	}{
		{0, 0},
		{1, 1},
		{1, 10},
		{1, 100},
		{1, 1000},
		{10, 1},
		{100, 1},
		{1000, 1},
		{4, 4}, // 341?
	}
	for _, c := range cases {
		b.Run(fmt.Sprintf("width=%d,depth=%d", c.width, c.depth), func(b *testing.B) {
			tracer := newTracer()

			rootSpan := tracer.StartSpan("a")
			ctx := opentracing.ContextWithSpan(context.Background(), rootSpan)

			root := &traceNode{}
			makeTree(root, c.width, c.depth)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				root.trace(ctx)
			}
		})
	}
}

type traceNode struct {
	children []*traceNode
}

func (t *traceNode) trace(ctx context.Context) {
	ctx, span := NewSpan(ctx, "hi")
	defer span.Finish()

	for _, child := range t.children {
		child.trace(ctx)
	}
}

type noopSpanPool struct{}

func (sp *noopSpanPool) get() *kSpan {
	return &kSpan{}
}

func (sp *noopSpanPool) put(s *kSpan) {
}
