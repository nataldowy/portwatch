package watch

import (
	"sync"
	"time"
)

// Debounce suppresses repeated events for the same key within a quiet period.
// Only the first event in a burst is allowed through; subsequent events within
// the wait duration are dropped. After the wait elapses with no new events the
// key is cleared so the next occurrence is allowed again.
type Debounce struct {
	mu      sync.Mutex
	wait    time.Duration
	timers  map[string]*time.Timer
	allowed map[string]bool
}

// NewDebounce creates a Debounce with the given quiet-period duration.
// A zero or negative wait defaults to 500 ms.
func NewDebounce(wait time.Duration) *Debounce {
	if wait <= 0 {
		wait = 500 * time.Millisecond
	}
	return &Debounce{
		wait:    wait,
		timers:  make(map[string]*time.Timer),
		allowed: make(map[string]bool),
	}
}

// Allow returns true the first time a key is seen within a burst and false for
// every subsequent call until the quiet period expires with no further events.
func (d *Debounce) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Reset the expiry timer on every call.
	if t, ok := d.timers[key]; ok {
		t.Reset(d.wait)
	} else {
		d.timers[key] = time.AfterFunc(d.wait, func() {
			d.mu.Lock()
			delete(d.timers, key)
			delete(d.allowed, key)
			d.mu.Unlock()
		})
	}

	if d.allowed[key] {
		return false
	}
	d.allowed[key] = true
	return true
}

// Reset clears all state, cancelling any pending timers.
func (d *Debounce) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for k, t := range d.timers {
		t.Stop()
		delete(d.timers, k)
	}
	for k := range d.allowed {
		delete(d.allowed, k)
	}
}
