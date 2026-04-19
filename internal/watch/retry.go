package watch

import (
	"sync"
	"time"
)

// Retry tracks per-key retry attempts with a fixed delay between tries.
type Retry struct {
	mu      sync.Mutex
	counts  map[string]int
	lastAt  map[string]time.Time
	delay   time.Duration
	maxTries int
	now     func() time.Time
}

// NewRetry creates a Retry that allows up to maxTries attempts per key,
// enforcing a minimum delay between consecutive tries.
func NewRetry(maxTries int, delay time.Duration) *Retry {
	return &Retry{
		counts:   make(map[string]int),
		lastAt:   make(map[string]time.Time),
		delay:    delay,
		maxTries: maxTries,
		now:      time.Now,
	}
}

// Allow returns true if the key may be retried (attempts remaining and delay elapsed).
func (r *Retry) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.counts[key] >= r.maxTries {
		return false
	}
	if last, ok := r.lastAt[key]; ok {
		if r.now().Sub(last) < r.delay {
			return false
		}
	}
	r.counts[key]++
	r.lastAt[key] = r.now()
	return true
}

// Remaining returns how many retry attempts are left for the key.
func (r *Retry) Remaining(key string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	remaining := r.maxTries - r.counts[key]
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears the retry state for the given key.
func (r *Retry) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.counts, key)
	delete(r.lastAt, key)
}
