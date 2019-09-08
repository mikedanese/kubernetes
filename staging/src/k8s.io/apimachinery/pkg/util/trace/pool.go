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
	"sync"
	"sync/atomic"
)

func newSpanPool(max int) *spanPool {
	return &spanPool{
		max: int32(max),
		pool: sync.Pool{
			New: func() interface{} {
				return &kSpan{}
			},
		},
	}
}

type pool interface {
	get() *kSpan
	put(*kSpan)
}

type spanPool struct {
	max  int32
	cur  int32
	pool sync.Pool
}

func (sp *spanPool) get() *kSpan {
	atomic.AddInt32(&sp.cur, 1)
	return sp.pool.Get().(*kSpan)
}

func (sp *spanPool) put(s *kSpan) {
	cur := atomic.LoadInt32(&sp.cur)
	if cur < sp.max {
		atomic.AddInt32(&sp.cur, -1)
		sp.pool.Put(s)
	}
}
