package report

import (
	"time"

	"portwatch/internal/history"
)

// Filter holds optional constraints for narrowing history entries.
type Filter struct {
	Kind  string    // "new", "closed", or "" for all
	Since time.Time // zero means no lower bound
	Until time.Time // zero means no upper bound
}

// Apply returns only the entries that match f.
func Apply(entries []history.Entry, f Filter) []history.Entry {
	out := make([]history.Entry, 0, len(entries))
	for _, e := range entries {
		if f.Kind != "" && e.Kind != f.Kind {
			continue
		}
		if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
			continue
		}
		if !f.Until.IsZero() && e.Timestamp.After(f.Until) {
			continue
		}
		out = append(out, e)
	}
	return out
}
