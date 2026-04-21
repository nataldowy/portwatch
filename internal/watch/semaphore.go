package watch

import "context"

// Semaphore limits the number of concurrent operations using a buffered channel.
type Semaphore struct {
	slots chan struct{}
}

// NewSemaphore creates a Semaphore with the given concurrency limit.
// If n < 1 it defaults to 1.
func NewSemaphore(n int) *Semaphore {
	if n < 1 {
		n = 1
	}
	return &Semaphore{slots: make(chan struct{}, n)}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns ctx.Err() if the context is done before a slot is acquired.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a previously acquired slot.
// It panics if called more times than Acquire has succeeded.
func (s *Semaphore) Release() {
	select {
	case <-s.slots:
	default:
		panic("semaphore: Release called without matching Acquire")
	}
}

// Cap returns the maximum number of concurrent holders.
func (s *Semaphore) Cap() int {
	return cap(s.slots)
}

// Len returns the number of currently held slots.
func (s *Semaphore) Len() int {
	return len(s.slots)
}

// Reset drains all held slots, returning the semaphore to an empty state.
func (s *Semaphore) Reset() {
	for {
		select {
		case <-s.slots:
		default:
			return
		}
	}
}
