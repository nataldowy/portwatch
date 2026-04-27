package watch

import (
	"sync"
	"time"
)

// Cascade propagates a signal across dependent keys: when a parent key fires,
// all registered children are suppressed for the given window duration.
type Cascade struct {
	mu       sync.Mutex
	children map[string][]string
	suppressed map[string]time.Time
	window   time.Duration
	now      func() time.Time
}

// NewCascade returns a Cascade with the given suppression window.
func NewCascade(window time.Duration) *Cascade {
	if window <= 0 {
		window = 5 * time.Second
	}
	return &Cascade{
		children:   make(map[string][]string),
		suppressed: make(map[string]time.Time),
		window:     window,
		now:        time.Now,
	}
}

// Register associates child keys with a parent key.
// When the parent fires, all children are suppressed.
func (c *Cascade) Register(parent string, children ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.children[parent] = append(c.children[parent], children...)
}

// Fire marks the parent as fired, suppressing all its children.
func (c *Cascade) Fire(parent string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiry := c.now().Add(c.window)
	for _, child := range c.children[parent] {
		c.suppressed[child] = expiry
	}
}

// Allow returns true if the given key is not currently suppressed.
func (c *Cascade) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiry, ok := c.suppressed[key]
	if !ok {
		return true
	}
	if c.now().After(expiry) {
		delete(c.suppressed, key)
		return true
	}
	return false
}

// Reset clears all suppression state and child registrations.
func (c *Cascade) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.children = make(map[string][]string)
	c.suppressed = make(map[string]time.Time)
}
