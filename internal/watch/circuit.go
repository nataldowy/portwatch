package watch

import (
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Circuit implements a per-key circuit breaker that opens after a threshold
// of consecutive failures and resets after a cooldown window.
type Circuit struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	failures  map[string]int
	openedAt  map[string]time.Time
	now       func() time.Time
}

func NewCircuit(threshold int, window time.Duration) *Circuit {
	return &Circuit{
		threshold: threshold,
		window:    window,
		failures:  make(map[string]int),
		openedAt:  make(map[string]time.Time),
		now:       time.Now,
	}
}

// Allow returns true if the key is allowed to proceed.
func (c *Circuit) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if t, open := c.openedAt[key]; open {
		if c.now().Sub(t) >= c.window {
			// transition to half-open: allow one probe
			delete(c.openedAt, key)
			c.failures[key] = 0
			return true
		}
		return false
	}
	return true
}

// RecordFailure records a failure for key and opens the circuit if threshold reached.
func (c *Circuit) RecordFailure(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.failures[key]++
	if c.failures[key] >= c.threshold {
		c.openedAt[key] = c.now()
	}
}

// RecordSuccess resets the failure count for key.
func (c *Circuit) RecordSuccess(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.failures, key)
	delete(c.openedAt, key)
}
