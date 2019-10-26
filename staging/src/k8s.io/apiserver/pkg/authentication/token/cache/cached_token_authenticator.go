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
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"io"
	"sync"
	"time"
	"unsafe"

	"k8s.io/apimachinery/pkg/util/cache"
	"k8s.io/apiserver/pkg/authentication/authenticator"
)

type cachedTokenAuthenticator struct {
	authenticator authenticator.Token

	cacheErrs  bool
	successTTL time.Duration
	failureTTL time.Duration

	cache *cache.Expiring

	// hashPool is a per authenticator pool of hash.Hash (to avoid allocations from building the Hash)
	// HMAC with SHA-256 and a random key is used to prevent precomputation and length extension attacks
	// It also mitigates hash map DOS attacks via collisions (the inputs are supplied by untrusted users)
	hashPool *sync.Pool
}

// New returns a token authenticator that caches the results of the specified authenticator. A ttl of 0 bypasses the cache.
func New(authenticator authenticator.Token, cacheErrs bool, successTTL, failureTTL time.Duration) authenticator.Token {
	randomCacheKey := make([]byte, 32)
	if _, err := rand.Read(randomCacheKey); err != nil {
		panic(err) // rand should never fail
	}

	return &cachedTokenAuthenticator{
		authenticator: authenticator,
		cacheErrs:     cacheErrs,
		successTTL:    successTTL,
		failureTTL:    failureTTL,
		cache:         cache.NewExpiring(),

		hashPool: &sync.Pool{
			New: func() interface{} {
				return hmac.New(sha256.New, randomCacheKey)
			},
		},
	}
}

// AuthenticateToken implements authenticator.Token
func (a *cachedTokenAuthenticator) AuthenticateToken(ctx context.Context, token string) (*authenticator.Response, bool, error) {
	auds, _ := authenticator.AudiencesFrom(ctx)

	key := keyFunc(a.hashPool, auds, token)
	if e, ok := a.cache.Get(key); ok {
		return e.(*authenticator.Response), true, nil
	}

	resp, ok, err := a.authenticator.AuthenticateToken(ctx, token)
	if err != nil {
		return resp, ok, err
	}

	if ok {
		a.cache.Set(key, resp, a.successTTL)
	}

	return resp, ok, err
}

// keyFunc generates a string key by hashing the inputs.
// This lowers the memory requirement of the cache and keeps tokens out of memory.
func keyFunc(hashPool *sync.Pool, auds []string, token string) string {
	h := hashPool.Get().(hash.Hash)

	h.Reset()

	// try to force stack allocation
	var a [4]byte
	b := a[:]

	writeLengthPrefixedString(h, b, token)
	// encode the length of audiences to avoid ambiguities
	writeLength(h, b, len(auds))
	for _, aud := range auds {
		writeLengthPrefixedString(h, b, aud)
	}

	key := toString(h.Sum(nil)) // skip base64 encoding to save an allocation

	hashPool.Put(h)

	return key
}

// writeLengthPrefixedString writes s with a length prefix to prevent ambiguities, i.e. "xy" + "z" == "x" + "yz"
// the length of b is assumed to be 4 (b is mutated by this function to store the length of s)
func writeLengthPrefixedString(w io.Writer, b []byte, s string) {
	writeLength(w, b, len(s))
	if _, err := w.Write(toBytes(s)); err != nil {
		panic(err) // Write() on hash never fails
	}
}

// writeLength encodes length into b and then writes it via the given writer
// the length of b is assumed to be 4
func writeLength(w io.Writer, b []byte, length int) {
	binary.BigEndian.PutUint32(b, uint32(length))
	if _, err := w.Write(b); err != nil {
		panic(err) // Write() on hash never fails
	}
}

// toBytes performs unholy acts to avoid allocations
func toBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// toString performs unholy acts to avoid allocations
func toString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
