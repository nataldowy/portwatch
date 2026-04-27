package watch

import (
	"testing"
	"time"
)

// TestShieldIntegratesWithPipeline verifies that a Shield used alongside the
// Pipeline correctly suppresses repeated alerts for the same port while
// allowing alerts for distinct ports through.
func TestShieldIntegratesWithPipeline(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	sh := NewShield(200*time.Millisecond, 0)

	port := testPort(9100, "new")

	// First event: pipeline allows, shield allows.
	if !pl.Allow(port) {
		t.Fatal("pipeline should allow first event")
	}
	if !sh.Allow(portKey(port)) {
		t.Fatal("shield should allow first event")
	}

	// Second event for same port: pipeline dedup blocks.
	if pl.Allow(port) {
		t.Fatal("pipeline dedup should block duplicate")
	}

	// Reset pipeline, but shield window still active.
	pl.Reset()
	if !pl.Allow(port) {
		t.Fatal("pipeline should allow after reset")
	}
	if sh.Allow(portKey(port)) {
		t.Fatal("shield should still block within window")
	}
}

// TestShieldWindowGrowsAcrossResets verifies that hit count and window growth
// survive pipeline resets (shield is independent state).
func TestShieldWindowGrowsAcrossResets(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	sh := NewShield(20*time.Millisecond, 10*time.Second)

	port := testPort(9200, "new")

	for round := 0; round < 3; round++ {
		pl.Reset()
		time.Sleep(25 * time.Millisecond) // wait for current shield window

		if !pl.Allow(port) {
			t.Fatalf("round %d: pipeline should allow", round)
		}
		if !sh.Allow(portKey(port)) {
			t.Fatalf("round %d: shield should allow after window", round)
		}
	}

	if sh.Hits(portKey(port)) != 3 {
		t.Fatalf("expected 3 hits after 3 rounds, got %d", sh.Hits(portKey(port)))
	}
}

func portKey(p Port) string {
	return itoa(p.Number) + "/" + p.Proto
}
