package watch

import (
	"sync"
	"testing"
)

func TestBufferConcurrentAdd(t *testing.T) {
	b := NewBuffer(128)
	var wg sync.WaitGroup
	keys := []string{"80", "443", "22", "8080"}
	for _, k := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			for i := 0; i < 25; i++ {
				b.Add(key)
			}
		}(k)
	}
	wg.Wait()
	entries := b.Flush()
	if len(entries) != len(keys) {
		t.Fatalf("expected %d entries, got %d", len(keys), len(entries))
	}
	for _, e := range entries {
		if e.Count != 25 {
			t.Errorf("key %s: expected count 25, got %d", e.Key, e.Count)
		}
	}
}

func TestBufferMaxUnderConcurrency(t *testing.T) {
	const max = 4
	b := NewBuffer(max)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			b.Add(string(rune('a' + n)))
		}(i)
	}
	wg.Wait()
	if b.Len() > max {
		t.Errorf("buffer exceeded max: got %d, want <= %d", b.Len(), max)
	}
}
