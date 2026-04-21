package watch

import (
	"sync"
	"time"
)

// Burst tracks short-term event spikes per key within a rolling window.
// It allows up to Cap events in Window duration; excess events are suppressed.
type Burst struct {
	mu     sync.Mutex
	cap    int
	window time.Duration
	buckets map[string]*burstBucket
}

type burstBucket struct {
	count  int
	reset  time.Time
}

// NewBurst creates a Burst filter allowing at most cap events per key per window.
// If cap < 1 it defaults to 1; if window <= 0 it defaults to 1 second.
func NewBurst(cap int, window time.Duration) *Burst {
	if cap < 1 {
		cap = 1
	}
	if window <= 0 {
		window = time.Second
	}
	return &Burst{
		cap:     cap,
		window:  window,
		buckets: make(map[string]*burstBucket),
	}
}

// Allow returns true if the event for key is within the burst cap.
func (b *Burst) Allow(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	bkt, ok := b.buckets[key]
	if !ok || now.After(bkt.reset) {
		b.buckets[key] = &burstBucket{count: 1, reset: now.Add(b.window)}
		return true
	}
	if bkt.count < b.cap {
		bkt.count++
		return true
	}
	return false
}

// Remaining returns how many events are still allowed for key in the current window.
func (b *Burst) Remaining(key string) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	bkt, ok := b.buckets[key]
	if !ok || now.After(bkt.reset) {
		return b.cap
	}
	rem := b.cap - bkt.count
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears the burst state for key.
func (b *Burst) Reset(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.buckets, key)
}
