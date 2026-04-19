package scanner

import (
	"testing"
	"time"
)

func makeSnapshot(ports ...int) Snapshot {
	var states []PortState
	for _, p := range ports {
		states = append(states, PortState{Port: p, Protocol: "tcp", Open: true})
	}
	return Snapshot{Ports: states, Timestamp: time.Now()}
}

func TestCompareDetectsNewPort(t *testing.T) {
	prev := makeSnapshot(80, 443)
	curr := makeSnapshot(80, 443, 8080)
	diff := Compare(prev, curr)
	if len(diff.Opened) != 1 || diff.Opened[0].Port != 8080 {
		t.Errorf("expected port 8080 to be opened, got %+v", diff.Opened)
	}
	if len(diff.Closed) != 0 {
		t.Errorf("expected no closed ports, got %+v", diff.Closed)
	}
}

func TestCompareDetectsClosedPort(t *testing.T) {
	prev := makeSnapshot(80, 443, 8080)
	curr := makeSnapshot(80, 443)
	diff := Compare(prev, curr)
	if len(diff.Closed) != 1 || diff.Closed[0].Port != 8080 {
		t.Errorf("expected port 8080 to be closed, got %+v", diff.Closed)
	}
	if len(diff.Opened) != 0 {
		t.Errorf("expected no opened ports, got %+v", diff.Opened)
	}
}

func TestCompareNoChanges(t *testing.T) {
	prev := makeSnapshot(80, 443)
	curr := makeSnapshot(80, 443)
	diff := Compare(prev, curr)
	if diff.HasChanges() {
		t.Error("expected no changes")
	}
}
