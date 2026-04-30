package watch

import (
	"sync"
	"time"
)

// Signal is a one-shot broadcast primitive. Once fired, all current and future
// waiters are immediately unblocked. It can be reset to rearm for the next
// event cycle.
type Signal struct {
	mu      sync.Mutex
	ch      chan struct{}
	fired   bool
	firedAt time.Time
}

// NewSignal returns an armed, unfired Signal.
func NewSignal() *Signal {
	return &Signal{ch: make(chan struct{})}
}

// Fire broadcasts to all waiters. Subsequent calls are no-ops until Reset.
func (s *Signal) Fire() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.fired {
		return
	}
	s.fired = true
	s.firedAt = time.Now()
	close(s.ch)
}

// Wait returns a channel that is closed when the signal fires.
func (s *Signal) Wait() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ch
}

// Fired reports whether the signal has been fired.
func (s *Signal) Fired() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.fired
}

// FiredAt returns the time the signal was fired, or the zero time if it has
// not yet been fired.
func (s *Signal) FiredAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.firedAt
}

// Reset rearms the signal so it can be fired again.
func (s *Signal) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ch = make(chan struct{})
	s.fired = false
	s.firedAt = time.Time{}
}
