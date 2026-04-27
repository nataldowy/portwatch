package watch

import (
	"testing"
)

// TestTallyTracksAlertFrequencyViaPipeline verifies that a Tally can be used
// alongside Pipeline to count how many times each port key passes the filter.
func TestTallyTracksAlertFrequencyViaPipeline(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	tl := NewTally()

	port := testPort("new", 9090)
	key := itoa(port.Port)

	// First occurrence should pass through pipeline.
	if !pl.Allow(port) {
		t.Fatal("expected first occurrence to be allowed")
	}
	tl.Record(key)

	// Second occurrence is a duplicate — pipeline blocks it.
	if pl.Allow(port) {
		t.Fatal("expected duplicate to be blocked")
	}
	// We only tally events that pass.

	e := tl.Get(key)
	if e == nil || e.Count != 1 {
		t.Fatalf("expected tally count 1, got %v", e)
	}
}

// TestTallyResetAfterPipelineResetAllowsRecount verifies that resetting both
// Tally and Pipeline restores counting from zero.
func TestTallyResetAfterPipelineResetAllowsRecount(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	tl := NewTally()

	port := testPort("new", 7070)
	key := itoa(port.Port)

	if !pl.Allow(port) {
		t.Fatal("first allow expected")
	}
	tl.Record(key)

	pl.Reset()
	tl.Reset()

	if !pl.Allow(port) {
		t.Fatal("expected allow after reset")
	}
	tl.Record(key)

	e := tl.Get(key)
	if e == nil || e.Count != 1 {
		t.Fatalf("expected count 1 after reset, got %v", e)
	}
}
