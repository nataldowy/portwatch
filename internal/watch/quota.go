package watch

import (
	"sync"
	"time"
)

// Quota enforces a maximum number of allowed events per key within a rolling
// time window, resetting automatically when the window expires.
type Quota struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	buckets map[string]*quotaBucket
}

type quotaBucket struct {
	count  int
	expiry time.Time
}

// NewQuota creates a Quota that allows up to max events per key within window.
// If max < 1 it defaults to 1; if window <= 0 it defaults to one minute.
func NewQuota(max int, window time.Duration) *Quota {
	if max < 1 {
		max = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &Quota{
		max:     max,
		window:  window,
		buckets: make(map[string]*quotaBucket),
	}
}

// Allow returns true and increments the counter if the key has not exhausted
// its quota for the current window. Returns false otherwise.
func (q *Quota) Allow(key string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	b, ok := q.buckets[key]
	if !ok || now.After(b.expiry) {
		q.buckets[key] = &quotaBucket{count: 1, expiry: now.Add(q.window)}
		return true
	}
	if b.count >= q.max {
		return false
	}
	b.count++
	return true
}

// Remaining returns how many more events the key may emit in its current window.
func (q *Quota) Remaining(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	b, ok := q.buckets[key]
	if !ok || time.Now().After(b.expiry) {
		return q.max
	}
	if q.max-b.count < 0 {
		return 0
	}
	return q.max - b.count
}

// Reset clears the quota state for the given key.
func (q *Quota) Reset(key string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.buckets, key)
}
