package watch

import (
	"sync"
	"time"
)

// Limiter enforces a token-bucket style global rate limit across all events.
type Limiter struct {
	mu       sync.Mutex
	rate     int
	window   time.Duration
	buckets  map[string][]time.Time
	nowFn    func() time.Time
}

// NewLimiter creates a Limiter allowing at most rate events per window globally per key.
func NewLimiter(rate int, window time.Duration) *Limiter {
	return &Limiter{
		rate:    rate,
		window:  window,
		buckets: make(map[string][]time.Time),
		nowFn:   time.Now,
	}
}

// Allow returns true if the key is within the allowed rate, recording the attempt.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.nowFn()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= l.rate {
		l.buckets[key] = filtered
		return false
	}

	l.buckets[key] = append(filtered, now)
	return true
}

// Reset clears the bucket for a given key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Remaining returns how many more events are allowed for key in the current window.
func (l *Limiter) Remaining(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.nowFn()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.buckets[key] {
		if t.After(cutoff) {
			count++
		}
	}
	r := l.rate - count
	if r < 0 {
		return 0
	}
	return r
}
