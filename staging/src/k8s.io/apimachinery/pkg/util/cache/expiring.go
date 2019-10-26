package cache

import (
	"sync"
	"time"
)

func NewExpiring() *Expiring {
	return &Expiring{
		now: time.Now,
		schedule: func(d time.Duration, f func()) (cancel func()) {
			timer := time.AfterFunc(d, f)
			return func() {
				timer.Stop()
			}
		},
		m: make(map[interface{}]entry),
	}
}

// Expiring is a map whose entries expire after a per-entry timeout. Its keys
// must all be of the same type, and that type must be a valid Go map key type.
type Expiring struct {
	now      func() time.Time
	schedule func(d time.Duration, f func()) (cancel func())

	// mu protects the below fields
	mu sync.RWMutex
	// m is the internal map that backs the cache.
	m map[interface{}]entry
	// gen is a generation number so that when setting multiple entries with the
	// same key, the expiration timer for earlier ones won't remove later ones.
	// It must be accessed atomicaly.
	gen uint64
}

type entry struct {
	val    interface{}
	expiry time.Time
	gen    uint64
	touch  func()
}

// Delete deletes an entry in the map.
func (c *Expiring) Delete(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.del(key, 0)
}

// Get looks up an entry in the cache.
func (c *Expiring) Get(key interface{}) (val interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.m[key]
	return e.val, ok
}

// Set sets a key/value/expiry entry in the map, overwriting any previous entry
// with the same key. The entry expires at the given expiry time, but it may be
// removed earlier by other Set or Clear calls.
func (c *Expiring) Set(key interface{}, val interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.del(key, 0)
	expiry := c.now().Add(ttl)

	c.gen++
	gen := c.gen

	c.m[key] = entry{
		val:    val,
		expiry: expiry,
		gen:    gen,
		touch: c.schedule(ttl, func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.del(key, gen)
		}),
	}
}

// del deletes the entry for the given key.
//
// If gen is non-zero, this only happens if the current entry has the same
// generation number as the one when the GC was scheudled. In most cases, the
// e.touch call (which cancels the GC) will prevent this, but strictly
// speaking, stopping a timer can race with the timer actually firing.
func (c *Expiring) del(key interface{}, gen uint64) {
	if c.m == nil {
		return
	}
	e, ok := c.m[key]
	if !ok {
		return
	}
	if gen != 0 && gen != e.gen {
		return
	}
	e.touch()
	delete(c.m, key)
}

// Len returns the number of items in the cache.
func (c *Expiring) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.m)
}
