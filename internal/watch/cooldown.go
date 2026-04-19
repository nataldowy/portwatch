package watch

import (
	"sync"
	"time"
)

// Cooldown tracks per-port suppression windows to avoid alert floods.
type Cooldown struct {
	mu       sync.Mutex
	window   time.Duration
	lastSeen map[string]time.Time
}

// NewCooldown creates a Cooldown with the given suppression window.
func NewCooldown(window time.Duration) *Cooldown {
	return &Cooldown{
		window:   window,
		lastSeen: make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// It records the current time for the key when returning true.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if last, ok := c.lastSeen[key]; ok && now.Sub(last) < c.window {
		return false
	}
	c.lastSeen[key] = now
	return true
}

// Reset clears the suppression record for a key.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.lastSeen, key)
}

// Len returns the number of tracked keys.
func (c *Cooldown) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.lastSeen)
}
