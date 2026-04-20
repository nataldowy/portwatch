package watch

import "sync"

// Counter tracks cumulative named counts across scan cycles.
// It is safe for concurrent use.
type Counter struct {
	mu     sync.Mutex
	counts map[string]int64
}

// NewCounter returns an initialised Counter.
func NewCounter() *Counter {
	return &Counter{
		counts: make(map[string]int64),
	}
}

// Inc increments the named counter by 1.
func (c *Counter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[key]++
}

// Add increments the named counter by delta. Negative values are ignored.
func (c *Counter) Add(key string, delta int64) {
	if delta <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[key] += delta
}

// Get returns the current value of the named counter.
func (c *Counter) Get(key string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[key]
}

// Snapshot returns a copy of all counters at a point in time.
func (c *Counter) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int64, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// Reset zeroes all counters.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[string]int64)
}
