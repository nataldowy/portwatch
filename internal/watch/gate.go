package watch

import (
	"sync"
	"time"
)

// Gate allows events through only when a minimum number of occurrences have
// been observed within a sliding window. This prevents alerting on transient
// single-occurrence spikes and requires confirmed repeated signals.
type Gate struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	events    map[string][]time.Time
	now       func() time.Time
}

// NewGate returns a Gate that allows an event through only after it has been
// seen at least threshold times within the given window duration.
func NewGate(threshold int, window time.Duration) *Gate {
	if threshold < 1 {
		threshold = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &Gate{
		threshold: threshold,
		window:    window,
		events:    make(map[string][]time.Time),
		now:       time.Now,
	}
}

// Allow records an occurrence for key and returns true if the threshold has
// been reached within the current window.
func (g *Gate) Allow(key string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := g.now()
	cutoff := now.Add(-g.window)

	// prune stale entries
	filtered := g.events[key][:0]
	for _, t := range g.events[key] {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	g.events[key] = filtered

	return len(filtered) >= g.threshold
}

// Count returns the number of recent occurrences for key within the window.
func (g *Gate) Count(key string) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := g.now()
	cutoff := now.Add(-g.window)
	count := 0
	for _, t := range g.events[key] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all recorded occurrences for key.
func (g *Gate) Reset(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.events, key)
}
