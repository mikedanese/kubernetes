// package lease implements common patterns for maintaining a lease on an
// expiring lock. Lease implementations should implement sync.Mutex interface
package lease

import (
	"time"

	"github.com/golang/glog"

	"k8s.io/kubernetes/pkg/util"
)

type LeaseLock interface {
	Events() <-chan LeaseEvent
}

type LeaseEvent string

const (
	LeaseEventAcquired LeaseEvent = "Acquired"
	LeaseEventLost                = "Lost"
)

func Lease(lock LeaseLock, onLease func(<-chan struct{}), onEviction func()) {
	var once bool
	lost := make(chan struct{}, 1)
	util.Until(func() {
		event := <-lock.Events()
		switch event {
		case LeaseEventAcquired:
			if once {
				onLease(lost)
			} else {
				glog.Error("LeaseEventAcquired sent down LeaseEvent channel more than once")
			}
		case LeaseEventLost:
			onEviction()
			lost <- struct{}{}
		default:
			panic("this is a programming error")
		}
	}, 5*time.Millisecond, lost)
}
