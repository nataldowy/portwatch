package watch

import (
	"sync"
	"time"
)

// Sampler records the last-seen timestamp for a key and reports whether
// enough time has elapsed since the previous sample.
type Sampler struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
	now      func() time.Time
}

// NewSampler returns a Sampler that allows an event through at most once per
// interval per key.
func NewSampler(interval time.Duration) *Sampler {
	return &Sampler{
		last:     make(map[string]time.Time),
		interval: interval,
		now:      time.Now,
	}
}

// Allow returns true when the key has not been seen within the interval.
func (s *Sampler) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	if t, ok := s.last[key]; ok && now.Sub(t) < s.interval {
		return false
	}
	s.last[key] = now
	return true
}

// Reset clears the recorded timestamp for a key.
func (s *Sampler) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.last, key)
}

// ResetAll clears all recorded timestamps.
func (s *Sampler) ResetAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = make(map[string]time.Time)
}

// LastSeen returns the time a key was last allowed through, and whether it
// has been seen at all.
func (s *Sampler) LastSeen(key string) (time.Time, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.last[key]
	return t, ok
}
