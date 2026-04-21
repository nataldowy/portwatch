package watch

import (
	"fmt"
	"sync"
	"testing"
)

func TestRegistryConcurrentRegister(t *testing.T) {
	r := NewRegistry()
	const workers = 20

	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("worker-%d", n)
			if err := r.Register(key, n); err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("unexpected registration error: %v", err)
	}

	if r.Len() != workers {
		t.Fatalf("expected %d entries, got %d", workers, r.Len())
	}
}

func TestRegistryConcurrentLookupAndUnregister(t *testing.T) {
	r := NewRegistry()
	const n = 50

	for i := 0; i < n; i++ {
		_ = r.Register(fmt.Sprintf("key-%d", i), i)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(2)
		go func(idx int) {
			defer wg.Done()
			r.Lookup(fmt.Sprintf("key-%d", idx))
		}(i)
		go func(idx int) {
			defer wg.Done()
			r.Unregister(fmt.Sprintf("key-%d", idx))
		}(i)
	}
	wg.Wait()
	// No panic or race — success.
}
