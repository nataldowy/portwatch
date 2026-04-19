package watch

import (
	"sync"
	"time"
)

// Throttle suppresses repeated alerts for the same port within a cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
}

// NewThrottle creates a Throttle with the given cooldown duration.
func NewThrottle(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// If allowed, the key's timestamp is updated.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.seen[key]; ok && now.Sub(last) < t.cooldown {
		return false
	}
	t.seen[key] = now
	return true
}

// Reset clears all throttle state.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.seen = make(map[string]time.Time)
}
