package watch

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// State holds the most recent snapshot and the time it was captured.
type State struct {
	mu       sync.RWMutex
	snapshot *scanner.Snapshot
	updated  time.Time
}

// Set stores a new snapshot.
func (s *State) Set(snap *scanner.Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshot = snap
	s.updated = time.Now()
}

// Get returns the current snapshot and its capture time.
func (s *State) Get() (*scanner.Snapshot, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshot, s.updated
}

// Age returns how long ago the snapshot was last updated.
func (s *State) Age() time.Duration {
	_, t := s.Get()
	if t.IsZero() {
		return 0
	}
	return time.Since(t)
}
