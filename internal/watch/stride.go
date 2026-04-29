package watch

import (
	"sync"
	"time"
)

// Stride tracks the rate of change between consecutive scans by measuring
// how many distinct keys fired within a sliding window. It signals when the
// rate exceeds a configured threshold, useful for detecting sudden bursts of
// port-change events across a short interval.
type Stride struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	events    map[string][]time.Time
	now       func() time.Time
}

// NewStride returns a Stride that fires when the number of unique keys seen
// within window exceeds threshold. A threshold < 1 is clamped to 1.
func NewStride(threshold int, window time.Duration) *Stride {
	if threshold < 1 {
		threshold = 1
	}
	if window <= 0 {
		window = 10 * time.Second
	}
	return &Stride{
		threshold: threshold,
		window:    window,
		events:    make(map[string][]time.Time),
		now:       time.Now,
	}
}

// Record registers an occurrence for key and reports whether the number of
// unique active keys now meets or exceeds the threshold.
func (s *Stride) Record(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	cutoff := now.Add(-s.window)

	// Append and prune stale timestamps for this key.
	s.events[key] = append(s.events[key], now)
	s.prune(key, cutoff)

	// Count keys that have at least one event within the window.
	active := 0
	for k := range s.events {
		s.prune(k, cutoff)
		if len(s.events[k]) > 0 {
			active++
		}
	}

	return active >= s.threshold
}

// ActiveKeys returns the number of keys that have recorded at least one event
// within the current window.
func (s *Stride) ActiveKeys() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := s.now().Add(-s.window)
	count := 0
	for k := range s.events {
		s.prune(k, cutoff)
		if len(s.events[k]) > 0 {
			count++
		}
	}
	return count
}

// Reset clears all recorded events.
func (s *Stride) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = make(map[string][]time.Time)
}

// prune removes timestamps older than cutoff for key. Must be called with mu held.
func (s *Stride) prune(key string, cutoff time.Time) {
	ts := s.events[key]
	i := 0
	for i < len(ts) && ts[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		s.events[key] = ts[i:]
	}
	if len(s.events[key]) == 0 {
		delete(s.events, key)
	}
}
