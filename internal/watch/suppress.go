package watch

import (
	"sync"
	"time"
)

// Suppress silences repeated alerts for a key until an explicit clear is
// called or the suppression window expires. Unlike Dedup, the window is
// extended on every hit, modelling a "stay quiet while the condition
// persists" pattern.
type Suppress struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]time.Time
	now     func() time.Time
}

// NewSuppress returns a Suppress with the given quiet window.
// A zero or negative window defaults to 30 seconds.
func NewSuppress(window time.Duration) *Suppress {
	if window <= 0 {
		window = 30 * time.Second
	}
	return &Suppress{
		window:  window,
		entries: make(map[string]time.Time),
		now:     time.Now,
	}
}

// Allow returns true if the key is not currently suppressed.
// Each call that returns false extends the suppression deadline by window.
func (s *Suppress) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	deadline, exists := s.entries[key]
	if exists && now.Before(deadline) {
		// Still suppressed — push the deadline further out.
		s.entries[key] = now.Add(s.window)
		return false
	}
	// First occurrence or window expired: record and allow.
	s.entries[key] = now.Add(s.window)
	return true
}

// Clear removes the suppression for key, allowing the next call to Allow
// to pass through immediately.
func (s *Suppress) Clear(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Reset removes all suppression state.
func (s *Suppress) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]time.Time)
}
