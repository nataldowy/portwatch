package watch

import (
	"sync"
	"time"
)

// Lease grants exclusive ownership of a key for a fixed duration.
// A second caller attempting to acquire the same key before the lease
// expires will be denied. Once the lease expires it can be re-acquired.
type Lease struct {
	mu      sync.Mutex
	leases  map[string]time.Time
	ttl     time.Duration
	now     func() time.Time
}

// NewLease returns a Lease with the given time-to-live per key.
// A zero or negative ttl is defaulted to 30 seconds.
func NewLease(ttl time.Duration) *Lease {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &Lease{
		leases: make(map[string]time.Time),
		ttl:    ttl,
		now:    time.Now,
	}
}

// Acquire attempts to obtain the lease for key.
// Returns true and the expiry time if the lease was granted.
// Returns false if the key is already leased and the lease has not expired.
func (l *Lease) Acquire(key string) (bool, time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if exp, ok := l.leases[key]; ok && now.Before(exp) {
		return false, exp
	}
	expiry := now.Add(l.ttl)
	l.leases[key] = expiry
	return true, expiry
}

// Release explicitly releases the lease for key before it expires.
// Calling Release on a key that has no active lease is a no-op.
func (l *Lease) Release(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.leases, key)
}

// Active returns true if key currently holds a valid (unexpired) lease.
func (l *Lease) Active(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	exp, ok := l.leases[key]
	return ok && l.now().Before(exp)
}

// Reset releases all active leases.
func (l *Lease) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.leases = make(map[string]time.Time)
}
