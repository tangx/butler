package butler

import (
	"context"
	"testing"
)

func BenchmarkButler(b *testing.B) {
	butler := Default()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	butler.WithOptions(WithContext(ctx)).Init()

	b.ResetTimer()
	go func() {
		for i := 0; i < 100000; i++ {
			// butler.AddJobs(func() { time.Sleep(1 * time.Millisecond) })
			butler.AddJobs(func() {})
		}
		cancel()
	}()

	butler.Work()
}
