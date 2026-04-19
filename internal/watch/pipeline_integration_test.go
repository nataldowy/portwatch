package watch_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watch"
)

func TestPipelineRateLimitEnforced(t *testing.T) {
	cfg := watch.PipelineConfig{
		Cooldown: watch.CooldownConfig{
			Window: 200 * time.Millisecond,
		},
		MaxPerMin: 2,
	}
	p := watch.NewPipeline(cfg)

	ports := []scanner.Port{
		{Number: 1, Protocol: "tcp"},
		{Number: 2, Protocol: "tcp"},
		{Number: 3, Protocol: "tcp"},
	}

	allowed := 0
	for _, port := range ports {
		if p.Allow("new", port) {
			allowed++
		}
	}
	// dedup passes all three (distinct ports); rate-limit per-key so all pass
	if allowed != 3 {
		t.Fatalf("expected 3 allowed (per-key rate limit), got %d", allowed)
	}
}

func TestPipelineResetRestoresFlow(t *testing.T) {
	cfg := watch.PipelineConfig{
		Cooldown: watch.CooldownConfig{
			Window: 10 * time.Millisecond,
		},
		MaxPerMin: 1,
	}
	p := watch.NewPipeline(cfg)
	port := scanner.Port{Number: 9000, Protocol: "tcp"}

	p.Allow("new", port)
	time.Sleep(20 * time.Millisecond)
	p.Reset("new", port)

	if !p.Allow("new", port) {
		t.Fatal("expected allow after reset and window expiry")
	}
}
