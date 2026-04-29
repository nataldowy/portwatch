package watch

import (
	"sync"
	"time"
)

// Stamp records the last-seen timestamp for a key and reports whether
// the gap since the previous occurrence exceeds a minimum interval.
// Unlike Fence, Stamp does not block subsequent calls — it only
// tracks timing and exposes metadata for downstream decisions.
type Stamp struct {
	mu       sync.Mutex
	interval time.Duration
	entries  map[string]time.Time
}

// NewStamp returns a Stamp that considers events "fresh" if they
// occurred within interval of the previous event for the same key.
// A zero or negative interval defaults to one second.
func NewStamp(interval time.Duration) *Stamp {
	if interval <= 0 {
		interval = time.Second
	}
	return &Stamp{
		interval: interval,
		entries:  make(map[string]time.Time),
	}
}

// Mark records the current time for key and returns true when the
// event is considered new (i.e. no prior record exists or the
// previous record is older than the configured interval).
func (s *Stamp) Mark(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	last, ok := s.entries[key]
	s.entries[key] = now
	if !ok {
		return true
	}
	return now.Sub(last) > s.interval
}

// LastSeen returns the time of the most recent Mark call for key
// and whether a record exists.
func (s *Stamp) LastSeen(key string) (time.Time, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.entries[key]
	return t, ok
}

// Age returns the duration since the last Mark for key.
// Returns zero and false if the key has never been marked.
func (s *Stamp) Age(key string) (time.Duration, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.entries[key]
	if !ok {
		return 0, false
	}
	return time.Since(t), true
}

// Reset clears all recorded timestamps.
func (s *Stamp) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]time.Time)
}
