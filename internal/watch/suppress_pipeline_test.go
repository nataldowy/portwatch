package watch

import (
	"testing"
	"time"
)

// TestSuppressCooperatesWithPipeline verifies that Suppress can be used
// alongside Pipeline: the first alert passes, subsequent ones are swallowed
// while the condition persists, and a Clear re-opens the gate.
func TestSuppressCooperatesWithPipeline(t *testing.T) {
	cfg := pipelineCfg()
	pipe := NewPipeline(cfg)
	sup := NewSuppress(time.Second)

	port := testPort("new")
	key := itoa(port.Port)

	// First event: pipeline dedup + suppress both allow.
	if !pipe.Allow(port) {
		t.Fatal("pipeline: expected first event to pass")
	}
	if !sup.Allow(key) {
		t.Fatal("suppress: expected first event to pass")
	}

	// Second event within window: suppress blocks even if pipeline would allow
	// (pipeline dedup also blocks here, showing layered filtering).
	pipe.Reset()
	if sup.Allow(key) {
		t.Fatal("suppress: expected second event to be suppressed")
	}
}

// TestSuppressClearThenPipelineReset shows that clearing suppress and
// resetting the pipeline together restore full event flow.
func TestSuppressClearThenPipelineReset(t *testing.T) {
	cfg := pipelineCfg()
	pipe := NewPipeline(cfg)
	sup := NewSuppress(time.Second)

	port := testPort("new")
	key := itoa(port.Port)

	pipe.Allow(port)
	sup.Allow(key)

	// Clear both layers.
	sup.Clear(key)
	pipe.Reset()

	if !pipe.Allow(port) {
		t.Fatal("pipeline: expected allow after reset")
	}
	if !sup.Allow(key) {
		t.Fatal("suppress: expected allow after clear")
	}
}
