package alert

import (
	"errors"
	"testing"
)

// stubNotifier records calls and optionally returns an error.
type stubNotifier struct {
	events []Event
	errOn bool
}

func (s *stubNotifier) Notify(e Event) error {
	s.events = append(s.events, e)
	if s.errOn {
		return errors.New("stub error")
	}
	return nil
}

func TestMultiNotifierFansOut(t *testing.T) {
	a, b := &stubNotifier{}, &stubNotifier{}
	mn := NewMultiNotifier(a, b)

	e := Event{Kind: "new", Port: makePort(8080)}
	if err := mn.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(a.events) != 1 || len(b.events) != 1 {
		t.Errorf("expected each notifier to receive 1 event, got a=%d b=%d", len(a.events), len(b.events))
	}
}

func TestMultiNotifierReturnsLastError(t *testing.T) {
	a := &stubNotifier{errOn: true}
	b := &stubNotifier{errOn: true}
	mn := NewMultiNotifier(a, b)

	e := Event{Kind: "new", Port: makePort(9090)}
	if err := mn.Notify(e); err == nil {
		t.Fatal("expected error, got nil")
	}
	// both notifiers should still have been called
	if len(a.events) != 1 || len(b.events) != 1 {
		t.Errorf("expected both notifiers called despite errors")
	}
}

func TestMultiNotifierAdd(t *testing.T) {
	mn := NewMultiNotifier()
	if mn.Len() != 0 {
		t.Fatalf("expected 0 notifiers")
	}
	mn.Add(&stubNotifier{})
	if mn.Len() != 1 {
		t.Fatalf("expected 1 notifier after Add")
	}
}

func TestMultiNotifierPartialError(t *testing.T) {
	ok := &stubNotifier{}
	bad := &stubNotifier{errOn: true}
	mn := NewMultiNotifier(bad, ok)

	e := Event{Kind: "closed", Port: makePort(443)}
	if err := mn.Notify(e); err == nil {
		t.Fatal("expected error from bad notifier")
	}
	// ok notifier must still receive the event
	if len(ok.events) != 1 {
		t.Errorf("ok notifier should have received event even after bad notifier failed")
	}
}
