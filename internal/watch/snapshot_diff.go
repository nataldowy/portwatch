package watch

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// DiffEvent holds a detected change between two snapshots.
type DiffEvent struct {
	Kind      string
	Port      scanner.Port
	DetectedAt time.Time
}

// SnapshotDiff compares two snapshots and returns a list of DiffEvents.
type SnapshotDiff struct{}

// NewSnapshotDiff creates a new SnapshotDiff.
func NewSnapshotDiff() *SnapshotDiff {
	return &SnapshotDiff{}
}

// Diff returns events for ports that appeared or disappeared.
func (sd *SnapshotDiff) Diff(prev, curr scanner.Snapshot) []DiffEvent {
	prevMap := portMap(prev.Ports)
	currMap := portMap(curr.Ports)
	var events []DiffEvent
	now := time.Now()

	for key, p := range currMap {
		if _, ok := prevMap[key]; !ok {
			events = append(events, DiffEvent{Kind: "new", Port: p, DetectedAt: now})
		}
	}
	for key, p := range prevMap {
		if _, ok := currMap[key]; !ok {
			events = append(events, DiffEvent{Kind: "closed", Port: p, DetectedAt: now})
		}
	}
	return events
}

func portMap(ports []scanner.Port) map[string]scanner.Port {
	m := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		key := p.Proto + ":" + itoa(p.Number)
		m[key] = p
	}
	return m
}
