package watch

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkThrottleGroupAllowSameKey(b *testing.B) {
	g := NewThrottleGroup(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Allow("port:8080:new")
	}
}

func BenchmarkThrottleGroupAllowUniqueKeys(b *testing.B) {
	g := NewThrottleGroup(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Allow(fmt.Sprintf("port:%d:new", i))
	}
}

func BenchmarkThrottleGroupActive(b *testing.B) {
	g := NewThrottleGroup(time.Second)
	for i := 0; i < 1000; i++ {
		g.Allow(fmt.Sprintf("port:%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Active()
	}
}
