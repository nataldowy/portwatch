package watch

import (
	"sync"
	"time"
)

// Fence enforces a minimum gap between consecutive allowed events per key.
// Unlike cooldown, it also tracks the total number of times a key was fenced
// (blocked) so callers can observe back-pressure.
type Fence struct {
	mu      sync.Mutex
	gap     time.Duration
	last    map[string]time.Time
	blocked map[string]int
	now     func() time.Time
}

// NewFence creates a Fence that enforces the given minimum gap between events.
// A zero or negative gap is replaced with 1 second.
func NewFence(gap time.Duration) *Fence {
	if gap <= 0 {
		gap = time.Second
	}
	return &Fence{
		gap:     gap,
		last:    make(map[string]time.Time),
		blocked: make(map[string]int),
		now:     time.Now,
	}
}

// Allow returns true if at least gap time has elapsed since the last allowed
// event for key. On denial the blocked counter for key is incremented.
func (f *Fence) Allow(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := f.now()
	if t, ok := f.last[key]; ok && now.Sub(t) < f.gap {
		f.blocked[key]++
		return false
	}
	f.last[key] = now
	return true
}

// Blocked returns how many times key has been denied since the last Reset.
func (f *Fence) Blocked(key string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.blocked[key]
}

// Reset clears all state for every key.
func (f *Fence) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.last = make(map[string]time.Time)
	f.blocked = make(map[string]int)
}
