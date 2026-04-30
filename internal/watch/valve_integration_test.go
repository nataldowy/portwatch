package watch

import (
	"sync"
	"testing"
	"time"
)

func TestValveConcurrentOpenClose(t *testing.T) {
	v := NewValve(time.Second)
	var wg sync.WaitGroup
	const workers = 50
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			v.Open("shared")
			_ = v.IsOpen("shared")
			v.Close("shared")
		}()
	}
	wg.Wait()
}

func TestValveConcurrentIndependentKeys(t *testing.T) {
	v := NewValve(time.Second)
	keys := []string{"alpha", "beta", "gamma", "delta"}
	var wg sync.WaitGroup
	wg.Add(len(keys) * 4)
	for _, k := range keys {
		k := k
		for j := 0; j < 4; j++ {
			go func() {
				defer wg.Done()
				v.Open(k)
				_ = v.IsOpen(k)
				v.Close(k)
				_ = v.IsOpen(k)
			}()
		}
	}
	wg.Wait()
}
