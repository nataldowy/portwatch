package watch

import (
	"sync"
	"time"
)

// Shield suppresses repeated alerts for a key until a minimum quiet period
// has elapsed since the last suppression was lifted. Unlike Cooldown, Shield
// tracks how many times a key was suppressed and increases the quiet window
// exponentially up to a configurable ceiling.
type Shield struct {
	mu      sync.Mutex
	window  time.Duration
	maxWin  time.Duration
	entries map[string]*shieldEntry
}

type shieldEntry struct {
	suppressedUntil time.Time
	hits            int
}

// NewShield creates a Shield with a base quiet window and a maximum window cap.
// If maxWindow is zero or less than window it defaults to 8× window.
func NewShield(window, maxWindow time.Duration) *Shield {
	if maxWindow <= window {
		maxWindow = window * 8
	}
	return &Shield{
		window:  window,
		maxWin:  maxWindow,
		entries: make(map[string]*shieldEntry),
	}
}

// Allow returns true when the key is not currently suppressed and records a
// new suppression window (doubled from the previous one, capped at maxWindow).
func (s *Shield) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	e, ok := s.entries[key]
	if ok && now.Before(e.suppressedUntil) {
		return false
	}

	nextWin := s.window
	if ok {
		nextWin = time.Duration(float64(s.window) * float64(uint(1)<<uint(e.hits)))
		if nextWin > s.maxWin {
			nextWin = s.maxWin
		}
	}

	hits := 0
	if ok {
		hits = e.hits + 1
	}
	s.entries[key] = &shieldEntry{
		suppressedUntil: now.Add(nextWin),
		hits:            hits,
	}
	return true
}

// Hits returns how many times Allow has returned true for the given key.
func (s *Shield) Hits(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.entries[key]; ok {
		return e.hits
	}
	return 0
}

// Reset clears all state for the given key.
func (s *Shield) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}
