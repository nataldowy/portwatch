package watch

import (
	"sync"
	"time"
)

// EventLogEntry records a single diff event for in-memory inspection.
type EventLogEntry struct {
	Kind  string
	Proto string
	Port  int
	At    time.Time
}

// EventLog is a thread-safe in-memory log of DiffEvents.
type EventLog struct {
	mu      sync.Mutex
	entries []EventLogEntry
	maxSize int
}

// NewEventLog creates an EventLog capped at maxSize entries (0 = unlimited).
func NewEventLog(maxSize int) *EventLog {
	return &EventLog{maxSize: maxSize}
}

// Append adds a DiffEvent to the log.
func (el *EventLog) Append(e DiffEvent) {
	el.mu.Lock()
	defer el.mu.Unlock()
	entry := EventLogEntry{
		Kind:  e.Kind,
		Proto: e.Port.Proto,
		Port:  e.Port.Number,
		At:    e.DetectedAt,
	}
	if el.maxSize > 0 && len(el.entries) >= el.maxSize {
		el.entries = el.entries[1:]
	}
	el.entries = append(el.entries, entry)
}

// All returns a copy of all log entries.
func (el *EventLog) All() []EventLogEntry {
	el.mu.Lock()
	defer el.mu.Unlock()
	out := make([]EventLogEntry, len(el.entries))
	copy(out, el.entries)
	return out
}

// Reset clears all entries.
func (el *EventLog) Reset() {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.entries = nil
}
