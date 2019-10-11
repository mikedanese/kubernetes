/*
Copyright 2017 The Kubernetes Authors.

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

package cache

import (
	"time"

	utilcache "k8s.io/apimachinery/pkg/util/cache"
	"k8s.io/apimachinery/pkg/util/clock"
)

type simpleCache struct {
	c *utilcache.Expiring
}

func newSimpleCache(_ int, _ clock.Clock) *simpleCache {
	return &simpleCache{
		c: utilcache.NewExpiring(),
	}
}

func (c *simpleCache) get(key string) (*cacheRecord, bool) {
	e, ok := c.c.Get(key)
	if !ok {
		return nil, false
	}
	return e.(*cacheRecord), true
}

func (c *simpleCache) set(key string, record *cacheRecord, ttl time.Duration) {
	c.c.Set(key, record, ttl)
}

func (c *simpleCache) remove(key string) {
	c.c.Delete(key)
}
