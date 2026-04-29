package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeRelayEvent(kind, port string) scanner.DiffEvent {
	return scanner.DiffEvent{
		Kind: kind,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
		At:   time.Now(),
	}
}

func TestRelayForwardsToSubscriber(t *testing.T) {
	r := NewRelay(nil)
	var got []scanner.DiffEvent
	r.Subscribe(func(ev scanner.DiffEvent) { got = append(got, ev) })

	ev := makeRelayEvent("new", "8080")
	r.Forward(ev)

	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	if got[0].Port.Number != "8080" {
		t.Errorf("unexpected port %s", got[0].Port.Number)
	}
}

func TestRelayFansOutToMultipleSubscribers(t *testing.T) {
	r := NewRelay(nil)
	count := 0
	for i := 0; i < 3; i++ {
		r.Subscribe(func(ev scanner.DiffEvent) { count++ })
	}

	r.Forward(makeRelayEvent("new", "443"))

	if count != 3 {
		t.Errorf("expected 3 calls, got %d", count)
	}
}

func TestRelayLenReflectsSubscribers(t *testing.T) {
	r := NewRelay(nil)
	if r.Len() != 0 {
		t.Fatal("expected 0 subscribers initially")
	}
	r.Subscribe(func(scanner.DiffEvent) {})
	r.Subscribe(func(scanner.DiffEvent) {})
	if r.Len() != 2 {
		t.Errorf("expected 2, got %d", r.Len())
	}
}

func TestRelayResetClearsHandlers(t *testing.T) {
	r := NewRelay(nil)
	r.Subscribe(func(scanner.DiffEvent) {})
	r.Reset()
	if r.Len() != 0 {
		t.Errorf("expected 0 after reset, got %d", r.Len())
	}

	var called bool
	r.Subscribe(func(scanner.DiffEvent) { called = true })
	r.Forward(makeRelayEvent("new", "22"))
	if !called {
		t.Error("handler registered after reset was not called")
	}
}

func TestRelayNoSubscribersNoError(t *testing.T) {
	r := NewRelay(nil)
	// should not panic
	r.Forward(makeRelayEvent("closed", "80"))
}
