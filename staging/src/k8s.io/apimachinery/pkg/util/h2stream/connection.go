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
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/httpstream"
)

func IsH2StreamRequest(req *http.Request) bool {
	if req.Proto != "HTTP/2.0" {
		return false
	}
	if !strings.EqualFold(req.Header.Get(httpstream.HeaderProtocol), "h2stream") {
		return false
	}
	return true
}

func NewServerConn(w http.ResponseWriter, req *http.Request) net.Conn {
	c := &conn{
		r: req.Body.Read,
		w: w.Write,
		c: func() error {
			return nil
		},
		closed: make(chan bool),
	}
	if f, ok := w.(http.Flusher); ok {
		c.f = f.Flush
	}
	return c
}

func NewClientConn(w io.Writer, resp *http.Response) net.Conn {
	return &conn{
		r:      resp.Body.Read,
		w:      w.Write,
		c:      resp.Body.Close,
		closed: make(chan bool),
	}
}

type conn struct {
	remoteAddr, localAddr net.Addr

	r func([]byte) (int, error)
	c func() error

	w func([]byte) (int, error)
	f func()

	closed chan bool
}

func (c *conn) Read(b []byte) (int, error) {
	return c.r(b)
}

func (c *conn) Write(b []byte) (int, error) {
	i, err := c.w(b)
	c.flush()
	return i, err
}

func (c *conn) Close() error {
	close(c.closed)
	return c.c()
}

func (c *conn) flush() {
	if c.f == nil {
		return
	}
	c.f()
}

func (c *conn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *conn) SetDeadline(t time.Time) error {
	return errors.New("unimplemented")
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return errors.New("unimplemented")
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return errors.New("unimplemented")
}
