package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(kind string, port int) DiffEvent {
	return DiffEvent{
		Kind:       kind,
		Port:       scanner.Port{Number: port, Proto: "tcp"},
		DetectedAt: time.Now(),
	}
}

func TestEventLogAppendAndAll(t *testing.T) {
	el := NewEventLog(0)
	el.Append(makeEvent("new", 80))
	el.Append(makeEvent("closed", 443))
	entries := el.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 80 || entries[0].Kind != "new" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
}

func TestEventLogCapsAtMaxSize(t *testing.T) {
	el := NewEventLog(3)
	for i := 0; i < 5; i++ {
		el.Append(makeEvent("new", i))
	}
	if len(el.All()) != 3 {
		t.Errorf("expected 3 entries after cap, got %d", len(el.All()))
	}
}

func TestEventLogOldestDroppedOnOverflow(t *testing.T) {
	el := NewEventLog(2)
	el.Append(makeEvent("new", 1))
	el.Append(makeEvent("new", 2))
	el.Append(makeEvent("new", 3))
	entries := el.All()
	if entries[0].Port != 2 {
		t.Errorf("expected oldest dropped, got port %d", entries[0].Port)
	}
}

func TestEventLogReset(t *testing.T) {
	el := NewEventLog(0)
	el.Append(makeEvent("new", 9000))
	el.Reset()
	if len(el.All()) != 0 {
		t.Error("expected empty log after reset")
	}
}
