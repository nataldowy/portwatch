package watch

import (
	"testing"
	"time"
)

// TestBudgetIntegratesWithPipeline verifies that a Budget used as a pre-filter
// ahead of a Pipeline correctly limits the total number of alerts dispatched.
func TestBudgetIntegratesWithPipeline(t *testing.T) {
	cfg := pipelineCfg()
	cfg.RateMax = 100 // effectively disable pipeline rate limiting
	p := NewPipeline(cfg)
	b := NewBudget(3, time.Minute)

	port := testPort("new")
	allowed := 0
	for i := 0; i < 6; i++ {
		if b.Allow("alerts", 1) {
			if p.Allow(port) {
				allowed++
				p.Reset()
			}
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 pipeline passes, got %d", allowed)
	}
}

// TestBudgetRemainingAfterPipelineDrop confirms that a pipeline rejection does
// not consume budget — budget is only spent when the caller decides to charge.
func TestBudgetRemainingAfterPipelineDrop(t *testing.T) {
	cfg := pipelineCfg()
	cfg.RateMax = 1
	p := NewPipeline(cfg)
	b := NewBudget(5, time.Minute)

	port := testPort("new")
	// First call passes both budget and pipeline.
	if b.Allow("k", 1) {
		p.Allow(port)
	}
	// Second call: pipeline blocks it; we choose NOT to charge budget.
	pipelineAllowed := p.Allow(port)
	if pipelineAllowed {
		b.Allow("k", 1)
	}
	// Budget should still have 4 remaining.
	if rem := b.Remaining("k"); rem != 4 {
		t.Fatalf("expected remaining 4, got %d", rem)
	}
}
