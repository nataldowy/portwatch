package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeSnapshotForDiff(ports []scanner.Port) scanner.Snapshot {
	return scanner.Snapshot{Ports: ports, At: time.Now()}
}

func TestSnapshotDiffDetectsNewPort(t *testing.T) {
	sd := NewSnapshotDiff()
	prev := makeSnapshotForDiff(nil)
	curr := makeSnapshotForDiff([]scanner.Port{{Number: 8080, Proto: "tcp"}})
	events := sd.Diff(prev, curr)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Kind != "new" {
		t.Errorf("expected kind=new, got %s", events[0].Kind)
	}
	if events[0].Port.Number != 8080 {
		t.Errorf("expected port 8080, got %d", events[0].Port.Number)
	}
}

func TestSnapshotDiffDetectsClosedPort(t *testing.T) {
	sd := NewSnapshotDiff()
	prev := makeSnapshotForDiff([]scanner.Port{{Number: 22, Proto: "tcp"}})
	curr := makeSnapshotForDiff(nil)
	events := sd.Diff(prev, curr)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Kind != "closed" {
		t.Errorf("expected kind=closed, got %s", events[0].Kind)
	}
}

func TestSnapshotDiffNoChanges(t *testing.T) {
	sd := NewSnapshotDiff()
	ports := []scanner.Port{{Number: 443, Proto: "tcp"}}
	prev := makeSnapshotForDiff(ports)
	curr := makeSnapshotForDiff(ports)
	events := sd.Diff(prev, curr)
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestSnapshotDiffTimestamp(t *testing.T) {
	sd := NewSnapshotDiff()
	before := time.Now()
	events := sd.Diff(
		makeSnapshotForDiff(nil),
		makeSnapshotForDiff([]scanner.Port{{Number: 9000, Proto: "udp"}}),
	)
	if len(events) == 0 {
		t.Fatal("expected event")
	}
	if events[0].DetectedAt.Before(before) {
		t.Error("DetectedAt should be >= test start")
	}
}
