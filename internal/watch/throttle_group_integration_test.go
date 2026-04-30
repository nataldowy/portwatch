package watch

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestThrottleGroupConcurrentAllow(t *testing.T) {
	g := NewThrottleGroup(time.Second)
	const goroutines = 50
	var wg sync.WaitGroup
	allowed := make([]bool, goroutines)
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			allowed[i] = g.Allow("shared-key")
		}()
	}
	wg.Wait()

	count := 0
	for _, a := range allowed {
		if a {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 allow for shared key under concurrency, got %d", count)
	}
}

func TestThrottleGroupConcurrentIndependentKeys(t *testing.T) {
	g := NewThrottleGroup(time.Second)
	const goroutines = 40
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("port:%d", i)
			if !g.Allow(key) {
				t.Errorf("expected first allow for unique key %s", key)
			}
		}()
	}
	wg.Wait()

	if got := g.Active(); got != goroutines {
		t.Fatalf("expected %d active keys, got %d", goroutines, got)
	}
}
