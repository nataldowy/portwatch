package watch

import (
	"sync"
	"time"
)

// Shed implements load-shedding: once the number of events recorded within a
// sliding window exceeds a high-water mark, further Allow calls return false
// until the load drops back below a low-water mark. This hysteresis prevents
// rapid toggling when traffic sits exactly at the threshold.
//
// Keys are tracked independently so that a noisy port does not affect others.
type Shed struct {
	mu        sync.Mutex
	high      int
	low       int
	window    time.Duration
	buckets   map[string]*shedBucket
}

type shedBucket struct {
	events   []time.Time
	shedding bool
}

// NewShed creates a Shed that starts shedding once high events are seen within
// window and stops shedding once the count drops to or below low.
// If high < 1 it is set to 1. If low < 1 or low >= high it is set to high-1.
// If window <= 0 it defaults to 10 seconds.
func NewShed(high, low int, window time.Duration) *Shed {
	if high < 1 {
		high = 1
	}
	if low < 1 || low >= high {
		low = high - 1
	}
	if window <= 0 {
		window = 10 * time.Second
	}
	return &Shed{
		high:    high,
		low:     low,
		window:  window,
		buckets: make(map[string]*shedBucket),
	}
}

// Allow records an event for key and returns true when the key is not currently
// being shed. The event is always recorded regardless of the return value so
// that the window count remains accurate.
func (s *Shed) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	b := s.bucket(key)
	b.events = pruneOld(b.events, now, s.window)
	b.events = append(b.events, now)

	count := len(b.events)

	if b.shedding {
		if count <= s.low {
			b.shedding = false
		}
	} else {
		if count >= s.high {
			b.shedding = true
		}
	}

	return !b.shedding
}

// Shedding reports whether key is currently in the load-shedding state.
func (s *Shed) Shedding(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.buckets[key]
	if !ok {
		return false
	}
	return b.shedding
}

// Count returns the number of events recorded for key within the current window.
func (s *Shed) Count(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.buckets[key]
	if !ok {
		return 0
	}
	now := time.Now()
	b.events = pruneOld(b.events, now, s.window)
	return len(b.events)
}

// Reset clears all state for key, including the shedding flag.
func (s *Shed) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.buckets, key)
}

func (s *Shed) bucket(key string) *shedBucket {
	if b, ok := s.buckets[key]; ok {
		return b
	}
	b := &shedBucket{}
	s.buckets[key] = b
	return b
}

// pruneOld removes events older than window from the front of the slice.
func pruneOld(events []time.Time, now time.Time, window time.Duration) []time.Time {
	cutoff := now.Add(-window)
	i := 0
	for i < len(events) && events[i].Before(cutoff) {
		i++
	}
	return events[i:]
}
