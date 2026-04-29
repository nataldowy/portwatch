package watch

import (
	"sync"
	"time"
)

// Epoch tracks a monotonically incrementing generation counter that advances
// whenever a significant boundary (e.g. a config reload or scan cycle reset)
// is crossed. Filters and pipelines can use the current epoch to discard
// stale state that belongs to a previous generation.
type Epoch struct {
	mu      sync.RWMutex
	current uint64
	updated time.Time
}

// NewEpoch returns an Epoch starting at generation 0.
func NewEpoch() *Epoch {
	return &Epoch{updated: time.Now()}
}

// Current returns the current epoch generation number.
func (e *Epoch) Current() uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.current
}

// Advance increments the epoch by one and records the time of the transition.
// It returns the new generation number.
func (e *Epoch) Advance() uint64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.current++
	e.updated = time.Now()
	return e.current
}

// Since returns how long ago the epoch last advanced.
func (e *Epoch) Since() time.Duration {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return time.Since(e.updated)
}

// Stale returns true when the supplied generation number is behind the
// current epoch, meaning the caller holds data from a previous generation.
func (e *Epoch) Stale(gen uint64) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return gen < e.current
}

// Reset sets the epoch back to zero and refreshes the updated timestamp.
func (e *Epoch) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.current = 0
	e.updated = time.Now()
}
