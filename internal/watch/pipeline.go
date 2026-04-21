package watch

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Pipeline chains dedup → cooldown → rate-limit filters for a single alert path.
// It is the primary guard between raw diff events and notifier dispatch.
type Pipeline struct {
	dedup     *Dedup
	cooldown  *Cooldown
	rateLimit *RateLimit
}

// PipelineCfg holds tunables for the Pipeline.
type PipelineCfg struct {
	Cooldown  time.Duration
	MaxPerMin int
}

// NewPipeline constructs a Pipeline with the provided configuration.
// Zero-value fields are replaced with sensible defaults.
func NewPipeline(cfg PipelineCfg) *Pipeline {
	cooldown := cfg.Cooldown
	if cooldown <= 0 {
		cooldown = 30 * time.Second
	}
	max := maxOrDefault(cfg.MaxPerMin, 10)
	return &Pipeline{
		dedup:     NewDedup(),
		cooldown:  NewCooldown(cooldown),
		rateLimit: NewRateLimit(max, time.Minute),
	}
}

// Allow returns true when the event passes all pipeline stages.
func (p *Pipeline) Allow(port scanner.Port, kind string) bool {
	key := itoa(port.Number) + "/" + kind
	if !p.dedup.Allow(port, kind) {
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

// Reset clears all pipeline state, useful between test runs or daemon restarts.
func (p *Pipeline) Reset() {
	p.dedup.Reset()
	p.cooldown.Reset()
	p.rateLimit.Reset()
}

func maxOrDefault(v, def int) int {
	if v > 0 {
		return v
	}
	return def
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
