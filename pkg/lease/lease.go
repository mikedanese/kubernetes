/*
 */

// package lease implements common patterns for maintaining a lease on an
// expiring lock. Lease implementations should implement sync.Mutex interface
package lease

import (
	"time"

	"k8s.io/kubernetes/pkg/util"
)

type LeaseLock interface {
	Events() chan LeaseEvent
}

type LeaseEvent string

const (
	LeaseEventAcquired LeaseEvent = "Acquired"
	LeaseEventRenewed             = "Renewed"
)

func Lease(lock LeaseLock, onLease func(chan struct{}), onEviction func(chan struct{})) {
	util.Until(func() {
		for {
			event := <-lock.Events()
			switch event {
			case LeaseEventAcquired:
			}
		}
	}, 5*time.Second, util.NeverStop)
}
