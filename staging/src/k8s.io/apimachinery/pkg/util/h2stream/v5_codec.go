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
	"io"

	"golang.org/x/net/http2"
)

func newV5Codec(c io.ReadWriter) codec {
	f := http2.NewFramer(c, c)
	f.SetReuseFrames()

	return &v5Codec{
		f: f,
	}
}

type v5Codec struct {
	f *http2.Framer
}

func (c *v5Codec) read() ([]byte, uint32, error) {
	for {
		f, err := c.f.ReadFrame()
		if err != nil {
			return nil, 0, err
		}
		switch f := f.(type) {
		case *http2.DataFrame:
			return f.Data(), f.Header().StreamID, nil
		}
	}
}

func (c *v5Codec) write(chid uint32, b []byte) error {
	return c.f.WriteData(chid, false, b)
}
