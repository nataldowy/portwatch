package watch

import (
	"sync"
	"time"
)

// RateLimit enforces a maximum number of alerts per window per key.
type RateLimit struct {
	mu      sync.Mutex
	window  time.Duration
	max     int
	buckets map[string]*rateBucket
}

type rateBucket struct {
	count int
	reset time.Time
}

// NewRateLimit creates a RateLimit allowing at most max events per window per key.
func NewRateLimit(window time.Duration, max int) *RateLimit {
	return &RateLimit{
		window:  window,
		max:     max,
		buckets: make(map[string]*rateBucket),
	}
}

// Allow returns true if the event for key is within the allowed rate.
func (r *RateLimit) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	b, ok := r.buckets[key]
	if !ok || now.After(b.reset) {
		r.buckets[key] = &rateBucket{count: 1, reset: now.Add(r.window)}
		return true
	}
	if b.count >= r.max {
		return false
	}
	b.count++
	return true
}

// Reset clears the bucket for the given key.
func (r *RateLimit) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.buckets, key)
}

// Remaining returns how many events are still allowed in the current window.
func (r *RateLimit) Remaining(key string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.buckets[key]
	if !ok || time.Now().After(b.reset) {
		return r.max
	}
	return r.max - b.count
}
