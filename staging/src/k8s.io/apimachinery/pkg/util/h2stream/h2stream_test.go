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
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"

	"golang.org/x/net/http2"

	certutil "k8s.io/client-go/util/cert"
)

type Testing interface {
	Errorf(format string, args ...interface{})
}

func newPipeStream(r io.ReadCloser, w io.WriteCloser) *Stream {
	return newStream(&conn{
		r: r.Read,
		w: w.Write,
		c: func() error {
			r.Close()
			w.Close()
			return nil
		}}, v5ProtocolVersion)
}

func setupBidi(t Testing) (*Stream, *Stream, func(context.Context), func()) {
	cr, sw := io.Pipe()
	sr, cw := io.Pipe()

	client := newPipeStream(cr, cw)

	server := newPipeStream(sr, sw)

	var wg sync.WaitGroup

	return client, server, func(ctx context.Context) {
			wg.Add(2)
			go func() {
				if err := client.Run(ctx); err != nil {
					t.Errorf("error running client stream: %v", err)
				}
				wg.Done()
			}()
			go func() {
				if err := server.Run(ctx); err != nil {
					t.Errorf("error running server stream: %v", err)
				}
				wg.Done()
			}()
		}, func() {
			wg.Wait()
		}
}

func pingPong(t Testing, ca, cb, sa, sb io.ReadWriteCloser) {
	if _, err := ca.Write([]byte("ca")); err != nil {
		t.Errorf("err: %v", err)
	}
	if _, err := sb.Write([]byte("sb")); err != nil {
		t.Errorf("err: %v", err)
	}

	cbuf := make([]byte, 2)
	cb.Read(cbuf)
	if got := string(cbuf); got != "sb" {
		t.Errorf("unexpected: %q", got)
	}

	sbuf := make([]byte, 2)
	sa.Read(sbuf)
	if got := string(sbuf); got != "ca" {
		t.Errorf("unexpected: %q", got)
	}
}

func newTestServer(t *testing.T, h http.Handler) *http.Server {
	certb, keyb, err := certutil.GenerateSelfSignedCertKey("testserver", nil, nil)
	if err != nil {
		t.Fatalf("err generating cert: %v", err)
	}
	cert, err := tls.X509KeyPair(certb, keyb)
	if err != nil {
		t.Fatalf("err parsing cert: %v", err)
	}

	ts := &http.Server{
		Addr:    "127.0.0.1:8999",
		Handler: h,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
		},
	}
	ln, err := net.Listen("tcp", ts.Addr)
	if err != nil {
		t.Fatalf("err starting test server: %v", err)
	}
	go ts.ServeTLS(ln, "", "")
	return ts
}

func doRequest(c *http.Client, method, url string) (*Stream, error) {
	r, w := io.Pipe()
	req, err := http.NewRequest(method, url, ioutil.NopCloser(r))
	if err != nil {
		return nil, err
	}
	req.Header.Set(protocolHeaderKey, protocolHeaderValue)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status was: %v", resp.Status)
	}
	return NewClientStream(w, resp), nil
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	r, w := io.Pipe()
	s := newPipeStream(r, w)

	wg.Add(1)
	go func() {
		s.Run(ctx)
		wg.Done()
	}()

	cancel()
	wg.Wait()
}

func TestReadWrite(t *testing.T) {
	client, server, run, _ := setupBidi(t)

	ca := client.Channel("a")
	cb := client.Channel("b")

	sa := server.Channel("a")
	sb := server.Channel("b")

	run(context.Background())

	pingPong(t, ca, cb, sa, sb)
}

func TestDropUnknown(t *testing.T) {
	client, server, run, wait := setupBidi(t)

	ca := client.Channel("a")
	cb := client.Channel("b")

	sa := server.Channel("a")
	sb := server.Channel("b")

	cu := client.Channel("cunknown")
	su := server.Channel("sunknown")

	ctx, cancel := context.WithCancel(context.Background())

	run(ctx)

	if _, err := cu.Write([]byte("cunknown")); err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, err := su.Write([]byte("sunknown")); err != nil {
		t.Fatalf("err: %v", err)
	}
	pingPong(t, ca, cb, sa, sb)

	cancel()
	wait()
}

func TestHTTP(t *testing.T) {
	mux := http.NewServeMux()
	var (
		sa, sb io.ReadWriteCloser
		err    error
	)
	mux.HandleFunc("/stream", func(w http.ResponseWriter, req *http.Request) {
		server, ok := Handle(w, req)
		if !ok {
			return
		}
		sa = server.Channel("a")
		sb = server.Channel("b")
		server.Run(req.Context())
	})
	ts := newTestServer(t, mux)
	defer ts.Shutdown(context.Background())

	cli := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := doRequest(cli, "GET", "https://"+ts.Addr+"/stream")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ca := client.Channel("a")
	cb := client.Channel("b")

	go client.Run(context.Background())

	pingPong(t, ca, cb, sa, sb)
}

func BenchmarkManyReadWrite(b *testing.B) {
	for _, numChannels := range []int{
		1,
		10,
		100,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprint(numChannels), func(b *testing.B) {
			client, server, run, _ := setupBidi(b)
			var wg sync.WaitGroup

			workers := []func(int){}

			for i := 0; i < numChannels; i++ {
				i := i
				wg.Add(1)

				ca := client.Channel(fmt.Sprintf("[%d]a", i))
				cb := client.Channel(fmt.Sprintf("[%d]b", i))

				sa := server.Channel(fmt.Sprintf("[%d]a", i))
				sb := server.Channel(fmt.Sprintf("[%d]b", i))

				workers = append(workers, func(numMessages int) {
					for j := 0; j < numMessages; j++ {
						pingPong(b, ca, cb, sa, sb)
					}
					wg.Done()
				})
			}

			run(context.Background())

			b.ResetTimer()
			for _, w := range workers {
				go w(b.N / numChannels)
			}

			wg.Wait()
		})
	}
}
