package watch

import (
	"sync"
	"time"
)

// Backoff tracks consecutive failures per key and returns an
// exponentially increasing wait duration before the next retry.
type Backoff struct {
	mu       sync.Mutex
	counts   map[string]int
	base     time.Duration
	max      time.Duration
}

// NewBackoff creates a Backoff with the given base and max durations.
func NewBackoff(base, max time.Duration) *Backoff {
	return &Backoff{
		counts: make(map[string]int),
		base:   base,
		max:    max,
	}
}

// Record increments the failure count for key and returns the delay
// that should be observed before the next attempt.
func (b *Backoff) Record(key string) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.counts[key]++
	return b.delay(b.counts[key])
}

// Reset clears the failure count for key.
func (b *Backoff) Reset(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.counts, key)
}

// Failures returns the current consecutive failure count for key.
func (b *Backoff) Failures(key string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.counts[key]
}

func (b *Backoff) delay(n int) time.Duration {
	d := b.base
	for i := 1; i < n; i++ {
		d *= 2
		if d >= b.max {
			return b.max
		}
	}
	if d > b.max {
		return b.max
	}
	return d
}
