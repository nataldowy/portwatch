package watch

import (
	"sync"
	"time"
)

// Valve is a flow-control primitive that opens and closes a named gate based
// on a sustained-open duration. Once opened it stays open until Close is
// called or the open window expires.
type Valve struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]valveEntry
}

type valveEntry struct {
	open    bool
	openAt  time.Time
	closedAt time.Time
}

// NewValve returns a Valve whose open window defaults to d.
// If d <= 0 it is clamped to one minute.
func NewValve(d time.Duration) *Valve {
	if d <= 0 {
		d = time.Minute
	}
	return &Valve{
		window:  d,
		entries: make(map[string]valveEntry),
	}
}

// Open marks key as open. Returns true if the key transitions from closed to
// open, false if it was already open.
func (v *Valve) Open(key string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	e := v.entries[key]
	if e.open {
		return false
	}
	v.entries[key] = valveEntry{open: true, openAt: time.Now()}
	return true
}

// Close marks key as closed. Returns true if the key was open.
func (v *Valve) Close(key string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	e, ok := v.entries[key]
	if !ok || !e.open {
		return false
	}
	e.open = false
	e.closedAt = time.Now()
	v.entries[key] = e
	return true
}

// IsOpen returns true when key is open and the open window has not expired.
func (v *Valve) IsOpen(key string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	e, ok := v.entries[key]
	if !ok || !e.open {
		return false
	}
	if time.Since(e.openAt) > v.window {
		e.open = false
		v.entries[key] = e
		return false
	}
	return true
}

// Reset clears all state.
func (v *Valve) Reset() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.entries = make(map[string]valveEntry)
}
