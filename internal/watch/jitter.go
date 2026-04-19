package watch

import (
	"math/rand"
	"sync"
	"time"
)

// Jitter adds random variance to a base duration to spread out bursts.
type Jitter struct {
	mu      sync.Mutex
	factor  float64 // fraction of base to vary, e.g. 0.2 = ±20%
	rng     *rand.Rand
}

// NewJitter creates a Jitter with the given factor (0 < factor <= 1).
func NewJitter(factor float64) *Jitter {
	if factor <= 0 || factor > 1 {
		factor = 0.1
	}
	return &Jitter{
		factor: factor,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Apply returns base ± (factor * base * random[0,1)).
func (j *Jitter) Apply(base time.Duration) time.Duration {
	j.mu.Lock()
	defer j.mu.Unlock()

	variance := float64(base) * j.factor
	delta := (j.rng.Float64()*2 - 1) * variance // range [-variance, +variance)
	result := time.Duration(float64(base) + delta)
	if result <= 0 {
		return base
	}
	return result
}

// Reset re-seeds the internal RNG.
func (j *Jitter) Reset() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}
