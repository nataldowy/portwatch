package alert

import (
	"testing"

	"portwatch/internal/scanner"
)

type captureNotifier struct {
	events []Event
}

func (c *captureNotifier) Notify(e Event) error {
	c.events = append(c.events, e)
	return nil
}

func TestDispatcherEmitsAlertForNewPort(t *testing.T) {
	cap := &captureNotifier{}
	d := NewDispatcher(cap)

	diff := scanner.Diff{
		New: []scanner.PortInfo{{Port: 9090, Proto: "tcp", PID: 1}},
	}
	d.Dispatch(diff)

	if len(cap.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(cap.events))
	}
	if cap.events[0].Level != LevelAlert {
		t.Errorf("expected ALERT level, got %s", cap.events[0].Level)
	}
}

func TestDispatcherEmitsWarnForClosedPort(t *testing.T) {
	cap := &captureNotifier{}
	d := NewDispatcher(cap)

	diff := scanner.Diff{
		Closed: []scanner.PortInfo{{Port: 22, Proto: "tcp", PID: 2}},
	}
	d.Dispatch(diff)

	if len(cap.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(cap.events))
	}
	if cap.events[0].Level != LevelWarn {
		t.Errorf("expected WARN level, got %s", cap.events[0].Level)
	}
}

func TestDispatcherNoEventsOnEmptyDiff(t *testing.T) {
	cap := &captureNotifier{}
	d := NewDispatcher(cap)
	d.Dispatch(scanner.Diff{})
	if len(cap.events) != 0 {
		t.Errorf("expected no events, got %d", len(cap.events))
	}
}
