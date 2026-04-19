package main

import (
	"context"
	"log"
	"time"

	"portwatch/internal/watch"
)

// supervisedRun wraps the provided run function with a Supervisor so that
// transient failures are retried before the process exits.
func supervisedRun(ctx context.Context, name string, maxRetries int, delay time.Duration, fn watch.Task) error {
	sup := watch.NewSupervisor(watch.SupervisorConfig{
		Name:       name,
		MaxRetries: maxRetries,
		Delay:      delay,
	})
	log.Printf("supervisor: starting %q (max retries=%d, delay=%s)", name, maxRetries, delay)
	return sup.Run(ctx, fn)
}
