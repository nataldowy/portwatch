package watch

import (
	"sync"
	"time"
)

// TallyEntry records the count and first/last seen times for a key.
type TallyEntry struct {
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Tally tracks occurrence counts per key with timestamps.
type Tally struct {
	mu      sync.Mutex
	entries map[string]*TallyEntry
}

// NewTally returns an initialised Tally.
func NewTally() *Tally {
	return &Tally{
		entries: make(map[string]*TallyEntry),
	}
}

// Record increments the count for key and updates timestamps.
func (t *Tally) Record(key string) int {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	if !ok {
		t.entries[key] = &TallyEntry{Count: 1, FirstSeen: now, LastSeen: now}
		return 1
	}
	e.Count++
	e.LastSeen = now
	return e.Count
}

// Get returns the TallyEntry for key, or nil if not present.
func (t *Tally) Get(key string) *TallyEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	if !ok {
		return nil
	}
	copy := *e
	return &copy
}

// Snapshot returns a copy of all entries.
func (t *Tally) Snapshot() map[string]TallyEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[string]TallyEntry, len(t.entries))
	for k, v := range t.entries {
		out[k] = *v
	}
	return out
}

// Reset clears all recorded entries.
func (t *Tally) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*TallyEntry)
}
