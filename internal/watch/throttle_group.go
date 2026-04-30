package watch

import (
	"sync"
	"time"
)

// ThrottleGroup manages per-key throttling across a shared group, allowing
// callers to suppress repeated events within a configurable window.
type ThrottleGroup struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]time.Time
}

// NewThrottleGroup returns a ThrottleGroup with the given suppression window.
// If window is <= 0 it defaults to 30 seconds.
func NewThrottleGroup(window time.Duration) *ThrottleGroup {
	if window <= 0 {
		window = 30 * time.Second
	}
	return &ThrottleGroup{
		window:  window,
		entries: make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the current window.
// On true it records the current time for the key.
func (g *ThrottleGroup) Allow(key string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := time.Now()
	if last, ok := g.entries[key]; ok && now.Sub(last) < g.window {
		return false
	}
	g.entries[key] = now
	return true
}

// Active returns the number of keys currently within their suppression window.
func (g *ThrottleGroup) Active() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := time.Now()
	count := 0
	for _, t := range g.entries {
		if now.Sub(t) < g.window {
			count++
		}
	}
	return count
}

// Reset removes all tracked entries.
func (g *ThrottleGroup) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.entries = make(map[string]time.Time)
}

// SetWindow updates the suppression window. Existing entries are unaffected.
func (g *ThrottleGroup) SetWindow(d time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if d > 0 {
		g.window = d
	}
}
