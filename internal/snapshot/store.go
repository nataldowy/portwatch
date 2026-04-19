package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"portwatch/internal/scanner"
)

// Record wraps a scanner snapshot with metadata for persistence.
type Record struct {
	SavedAt  time.Time        `json:"saved_at"`
	Snapshot scanner.Snapshot `json:"snapshot"`
}

// Store persists and loads port snapshots to/from disk.
type Store struct {
	path string
}

// NewStore creates a Store that reads/writes to the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the snapshot to disk, creating parent directories as needed.
func (s *Store) Save(snap scanner.Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	rec := Record{SavedAt: time.Now(), Snapshot: snap}
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(rec)
}

// Load reads the most recently saved snapshot from disk.
// Returns (zero Record, false, nil) when no file exists yet.
func (s *Store) Load() (Record, bool, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return Record{}, false, nil
	}
	if err != nil {
		return Record{}, false, err
	}
	defer f.Close()
	var rec Record
	if err := json.NewDecoder(f).Decode(&rec); err != nil {
		return Record{}, false, err
	}
	return rec, true, nil
}
