package lease

import (
	"time"

	etcdstorage "k8s.io/kubernetes/pkg/storage/etcd"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

type etcdLease struct {
	holderID  string
	ttl       time.Duration
	lastLease time.Time
	key       string
	client    *etcd.Client
}

// runs the election loop. never returns.
func (l *etcdLease) leaseAndUpdateLoop() {
	for {
		leaseHeld, err := l.acquireOrRenewLease()
		if err != nil {
			glog.Errorf("Error in master election: %v", err)
			if time.Now().Sub(l.lastLease) < l.ttl {
				continue
			}
			// Our lease has expired due to our own accounting, pro-actively give it
			// up, even if we couldn't contact etcd.
			glog.Infof("Too much time has elapsed, giving up lease.")
			leaseHeld = false
		}
		if err := l.update(leaseHeld); err != nil {
			glog.Errorf("Error updating files: %v", err)
		}
		//time.Sleep(l.sleep)
	}
}

// acquireOrRenewLease either races to acquire a new master lease, or update the existing master's lease
// returns true if we have the lease, and an error if one occurs.
// TODO: use the master election utility once it is merged in.
func (l *etcdLease) acquireOrRenewLease() (bool, error) {
	result, err := l.client.Get(l.key, false, false)
	if err != nil {
		if etcdstorage.IsEtcdNotFound(err) {
			// there is no current lease, try to acquire lease, create will fail if the key already exists
			_, err := l.client.Create(l.key, l.holderID, uint64(l.ttl))
			if err != nil {
				return false, err
			}
			l.lastLease = time.Now()
			return true, nil
		}
		return false, err
	}
	if result.Node.Value == l.holderID {
		glog.Infof("key already exists, we are the master (%s)", result.Node.Value)
		// we extend our lease @ 1/2 of the existing TTL, this ensures the master doesn't flap around
		if result.Node.Expiration.Sub(time.Now()) < l.ttl/2 {
			_, err := l.client.CompareAndSwap(l.key, l.holderID, uint64(l.ttl), l.holderID, result.Node.ModifiedIndex)
			if err != nil {
				return false, err
			}
		}
		l.lastLease = time.Now()
		return true, nil
	}
	glog.Infof("key already exists, the master is %s, sleeping.", result.Node.Value)
	return false, nil
}

func (l *etcdLease) update(leaseHeld bool) error {
	return nil
}
