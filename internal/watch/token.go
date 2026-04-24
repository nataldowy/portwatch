package watch

import (
	"sync"
	"time"
)

// Token implements a token-bucket rate limiter that replenishes tokens
// at a fixed rate up to a configurable capacity.
type Token struct {
	mu       sync.Mutex
	tokens   map[string]float64
	last     map[string]time.Time
	rate     float64 // tokens per second
	capacity float64
	now      func() time.Time
}

// NewToken returns a Token bucket limiter with the given replenish rate
// (tokens/sec) and burst capacity. Invalid values default to 1.
func NewToken(rate, capacity float64) *Token {
	if rate <= 0 {
		rate = 1
	}
	if capacity <= 0 {
		capacity = 1
	}
	return &Token{
		tokens:   make(map[string]float64),
		last:     make(map[string]time.Time),
		rate:     rate,
		capacity: capacity,
		now:      time.Now,
	}
}

// Allow returns true and consumes one token for the given key if a token
// is available. Returns false if the bucket is empty.
func (t *Token) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[key]; ok {
		elapsed := now.Sub(last).Seconds()
		t.tokens[key] += elapsed * t.rate
		if t.tokens[key] > t.capacity {
			t.tokens[key] = t.capacity
		}
	} else {
		t.tokens[key] = t.capacity
	}
	t.last[key] = now

	if t.tokens[key] >= 1 {
		t.tokens[key]--
		return true
	}
	return false
}

// Remaining returns the current token count for the given key (floored).
func (t *Token) Remaining(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	v, ok := t.tokens[key]
	if !ok {
		return int(t.capacity)
	}
	if v < 0 {
		return 0
	}
	return int(v)
}

// Reset clears all state for the given key.
func (t *Token) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.tokens, key)
	delete(t.last, key)
}
