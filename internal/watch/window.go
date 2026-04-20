package watch

import (
	"sync"
	"time"
)

// Window tracks event counts within a rolling time window per key.
type Window struct {
	mu       sync.Mutex
	buckets  map[string][]time.Time
	duration time.Duration
	max      int
}

// NewWindow creates a Window that allows at most max events per key within d.
func NewWindow(d time.Duration, max int) *Window {
	if d <= 0 {
		d = time.Minute
	}
	if max <= 0 {
		max = 10
	}
	return &Window{
		buckets:  make(map[string][]time.Time),
		duration: d,
		max:      max,
	}
}

// Allow returns true if the event for key is within the allowed rate.
func (w *Window) Allow(key string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	w.evict(key, now)
	if len(w.buckets[key]) >= w.max {
		return false
	}
	w.buckets[key] = append(w.buckets[key], now)
	return true
}

// Remaining returns how many more events are allowed for key right now.
func (w *Window) Remaining(key string) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(key, time.Now())
	r := w.max - len(w.buckets[key])
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all recorded events for key.
func (w *Window) Reset(key string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.buckets, key)
}

func (w *Window) evict(key string, now time.Time) {
	cutoff := now.Add(-w.duration)
	times := w.buckets[key]
	i := 0
	for i < len(times) && times[i].Before(cutoff) {
		i++
	}
	w.buckets[key] = times[i:]
}
