package watch

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkValveOpenIsOpenClose(b *testing.B) {
	v := NewValve(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Open("k")
		_ = v.IsOpen("k")
		v.Close("k")
	}
}

func BenchmarkValveUniqueKeys(b *testing.B) {
	v := NewValve(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("port-%d", i)
		v.Open(key)
		_ = v.IsOpen(key)
	}
}

func BenchmarkValveIsOpenHotPath(b *testing.B) {
	v := NewValve(time.Hour)
	v.Open("hot")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.IsOpen("hot")
	}
}
