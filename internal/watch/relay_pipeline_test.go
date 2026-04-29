package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// TestRelayFeedsIntoDedup verifies that a Relay can be composed with Dedup
// so that duplicate events are suppressed at the downstream handler.
func TestRelayFeedsIntoDedup(t *testing.T) {
	dedup := NewDedup()
	relay := NewRelay(nil)

	var passed int
	relay.Subscribe(func(ev scanner.DiffEvent) {
		if dedup.Allow(ev.Port.Number, ev.Kind) {
			passed++
		}
	})

	ev := scanner.DiffEvent{
		Kind: "new",
		Port: scanner.Port{Number: "3000", Protocol: "tcp"},
		At:   time.Now(),
	}

	relay.Forward(ev)
	relay.Forward(ev) // duplicate — should be blocked by dedup

	if passed != 1 {
		t.Errorf("expected 1 event through dedup, got %d", passed)
	}
}

// TestRelayResetThenResubscribe ensures that after Reset, newly registered
// handlers receive events while old ones do not.
func TestRelayResetThenResubscribe(t *testing.T) {
	relay := NewRelay(nil)

	oldCalled := false
	relay.Subscribe(func(scanner.DiffEvent) { oldCalled = true })
	relay.Reset()

	newCalled := false
	relay.Subscribe(func(scanner.DiffEvent) { newCalled = true })

	relay.Forward(scanner.DiffEvent{
		Kind: "closed",
		Port: scanner.Port{Number: "22", Protocol: "tcp"},
		At:   time.Now(),
	})

	if oldCalled {
		t.Error("old handler should not be called after Reset")
	}
	if !newCalled {
		t.Error("new handler should have been called")
	}
}
