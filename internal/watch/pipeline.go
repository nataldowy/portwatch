package watch

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Filter is a function that decides whether an alert should proceed.
type Filter func(kind, key string) bool

// Pipeline chains dedup, cooldown and rate-limit filters before dispatching.
type Pipeline struct {
	dedup     *Dedup
	cooldown  *Cooldown
	rateLimit *RateLimit
}

// PipelineConfig holds tunables for the Pipeline.
type PipelineConfig struct {
	Cooldown  CooldownConfig
	MaxPerMin int
}

// NewPipeline constructs a Pipeline from the given config.
func NewPipeline(cfg PipelineConfig) *Pipeline {
	return &Pipeline{
		dedup:     NewDedup(),
		cooldown:  NewCooldown(cfg.Cooldown),
		rateLimit: NewRateLimit(cfg.Cooldown.Window, maxOrDefault(cfg.MaxPerMin, 10)),
	}
}

// Allow returns true when the event passes all filters.
func (p *Pipeline) Allow(kind string, port scanner.Port) bool {
	key := fmt.Sprintf("%s:%d", kind, port.Number)
	if !p.dedup.Allow(kind, port) {
		return false
	}
	if !p.cooldown.Allow(key) {
		return false
	}
	if !p.rateLimit.Allow(key) {
		return false
	}
	return true
}

// Reset clears all filter state for a given kind+port combination.
func (p *Pipeline) Reset(kind string, port scanner.Port) {
	key := fmt.Sprintf("%s:%d", kind, port.Number)
	p.dedup.Reset()
	p.cooldown.Reset(key)
	p.rateLimit.Reset(key)
}

func maxOrDefault(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}
