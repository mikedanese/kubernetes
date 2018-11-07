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
	"bytes"
	"io"
	"sync"
)

func newPipe() io.ReadWriteCloser {
	var (
		mu sync.Mutex
		p  pipe
	)
	p.Cond.L = &mu
	return &p
}

// pipe implements a pipe that buffers writes and blocks reads if the buffer is
// empty until the next write, or until the buffer is closed.
type pipe struct {
	sync.Cond

	closed bool
	buf    bytes.Buffer
}

func (p *pipe) Read(b []byte) (int, error) {
	p.L.Lock()
	defer p.L.Unlock()

	for {
		if p.buf.Len() > 0 {
			return p.buf.Read(b)
		}
		if p.closed {
			return 0, io.EOF
		}
		p.Wait()
	}
}

func (p *pipe) Write(b []byte) (int, error) {
	p.L.Lock()
	defer p.L.Unlock()
	defer p.Signal()

	if p.closed {
		return 0, io.ErrClosedPipe
	}

	return p.buf.Write(b)
}

func (p *pipe) Close() error {
	p.L.Lock()
	defer p.L.Unlock()
	defer p.Broadcast()

	p.closed = true
	return nil
}
