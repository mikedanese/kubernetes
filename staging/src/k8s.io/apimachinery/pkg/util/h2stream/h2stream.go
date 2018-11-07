/*
Copyright 2018 The Kubernetes Authors.

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

package h2stream

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	protocolHeaderKey     = "protocol"
	protocolHeaderValue   = "h2stream"
	protocolVersionHeader = "X-Stream-Protocol-Version"

	v5ProtocolVersion = "v5.channel.k8s.io"
)

func IsH2StreamRequest(req *http.Request) bool {
	if req.Proto != "HTTP/2.0" {
		return false
	}
	if !strings.EqualFold(req.Header.Get(protocolHeaderKey), protocolHeaderValue) {
		return false
	}
	return true
}

func Handle(w http.ResponseWriter, req *http.Request) (*Stream, bool) {
	if !IsH2StreamRequest(req) {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("not a h2 stream"))
		return nil, false
	}
	w.Header().Set(protocolVersionHeader, v5ProtocolVersion)
	w.WriteHeader(http.StatusOK)
	return NewServerStream(w, req), true
}

func NewServerConn(w http.ResponseWriter, req *http.Request) *conn {
	c := &conn{
		r: req.Body.Read,
		w: w.Write,
		c: func() error {
			return nil
		},
	}
	if f, ok := w.(http.Flusher); ok {
		c.f = f.Flush
	}
	return c
}

func NewServerStream(w http.ResponseWriter, req *http.Request) *Stream {
	return newStream(NewServerConn(w, req), w.Header().Get(protocolVersionHeader))
}

func NewClientConn(w io.Writer, resp *http.Response) *conn {
	return &conn{
		r: resp.Body.Read,
		w: w.Write,
		c: resp.Body.Close,
	}
}

func NewClientStream(w io.Writer, resp *http.Response) *Stream {
	return newStream(NewClientConn(w, resp), resp.Header.Get(protocolVersionHeader))
}

func newStream(c *conn, selectedProtocol string) *Stream {
	return &Stream{
		needsFlush: make(chan struct{}, 1),
		conn:       c,
		codec:      selectCodec(c, selectedProtocol),
		chmap:      make(map[uint32]*channel),
	}
}

type Stream struct {
	needsFlush chan struct{}
	// guards writes to codec or underlying conn, including flushes
	wmu   sync.Mutex
	conn  *conn
	codec codec

	chmap map[uint32]*channel
}

func (s *Stream) Channel(name string) io.ReadWriteCloser {
	chash := fnv.New32()
	chash.Write([]byte(name))
	chid := chash.Sum32() &^ (1 << 31)

	ch := &channel{
		id:  chid,
		rwc: newPipe(),
		write: func(b []byte) (int, error) {
			s.wmu.Lock()
			defer s.wmu.Unlock()

			select {
			case s.needsFlush <- struct{}{}:
			default:
			}

			if err := s.codec.write(chid, b); err != nil {
				return 0, err
			}

			return len(b), nil
		},
	}

	s.chmap[chid] = ch
	return ch
}

func (s *Stream) Run(ctx context.Context) error {
	go func() {
		// flush once on start to get the headers written and the channel going.
		select {
		case s.needsFlush <- struct{}{}:
		default:
		}

		for {
			select {
			case <-ctx.Done():
				s.conn.close()
				return
			case <-s.needsFlush:
				s.wmu.Lock()
				s.conn.flush()
				s.wmu.Unlock()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	for {
		if err := s.recieveData(); err != nil {
			switch err {
			case io.EOF, io.ErrClosedPipe:
				return nil
			default:
				return err
			}
		}
	}
}

func (s *Stream) recieveData() error {
	b, chid, err := s.codec.read()
	if err != nil {
		return err
	}
	ch, ok := s.chmap[chid]
	if !ok {
		// consider logging a warning
		return nil
	}
	ch.rwc.Write(b)
	return nil
}

func (s *Stream) Close() error {
	s.conn.flush()
	return nil
}

func selectCodec(c *conn, selectedProtocol string) codec {
	switch selectedProtocol {
	case v5ProtocolVersion:
		return newV5Codec(c)
	default:
		panic("unknown protocol: " + selectedProtocol)
	}
}

type codec interface {
	read() ([]byte, uint32, error)
	write(chid uint32, data []byte) error
}

type conn struct {
	r func([]byte) (int, error)
	c func() error

	w func([]byte) (int, error)
	f func()
}

func (c *conn) Read(b []byte) (int, error) {
	return c.r(b)
}

func (c *conn) Write(b []byte) (int, error) {
	return c.w(b)
}

func (c *conn) close() error {
	return c.c()
}

func (c *conn) flush() {
	if c.f == nil {
		return
	}
	c.f()
}

type channel struct {
	id    uint32
	rwc   io.ReadWriteCloser
	write func([]byte) (int, error)
}

func (c *channel) Read(b []byte) (int, error) {
	return c.rwc.Read(b)
}

func (c *channel) Write(b []byte) (int, error) {
	return c.write(b)
}

func (c *channel) Close() error { return c.rwc.Close() }
