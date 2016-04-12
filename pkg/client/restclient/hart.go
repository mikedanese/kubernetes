/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package restclient

import (
	"math/rand"
	"net/http"
	"sync"
)

type hostProvider interface {
	get() string
	reportFailure(host string)
	reset()
}

type highAvailibilityRoundTripper struct {
	once     sync.Once
	delegate http.RoundTripper
	provider hostProvider
}

func (hart *highAvailibilityRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	host := hart.provider.get()
	resp, err := hart.delegate.RoundTrip(req)
	if err != nil {
		hart.provider.reportFailure(host)
	} else {
		hart.provider.reset()
	}
	return resp, err
}

type slightlyStickyProvider struct {
	hosts     []string
	threshold int
	rand      rand.Rand

	// bookkeeping
	sync.RWMutex
	fails int
	cur   int
}

func (s *slightlyStickyProvider) get() string {
	s.RLock()
	defer s.RUnlock()
	return s.hosts[s.cur]
}

func (s *slightlyStickyProvider) reportFailure(host string) {
	s.Lock()
	defer s.Unlock()

	if s.hosts[s.cur] != host {
		return
	}
	s.fails += 1
	if s.fails > s.threshold {
		s.cur = s.rand.Intn(len(s.hosts))
		s.fails = 0
	}
}

func (s *slightlyStickyProvider) reset() {
	s.Lock()
	defer s.Unlock()
	s.fails = 0
}
