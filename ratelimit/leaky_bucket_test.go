package ratelimit

import (
	"context"
	"math"
	"testing"
	"time"
)

func TestLeakyBucketRateLimiter(t *testing.T) {
	t.Run("Test1", func(t *testing.T) {
		done := make(chan bool, 20)

		rateLimiter := NewLeakyBucketRateLimiter(context.TODO(), 10, 5)
		start_t := time.Now()

		for i := 0; i < 20; i++ {
			go func() {
				rateLimiter.WaitForToken()
				done <- true
			}()
		}

		for i := 0; i < 20; i++ {
			<-done
		}

		elapsed := time.Since(start_t).Milliseconds()
		if math.Abs(float64(elapsed)-2000.0) > 30 {
			t.Errorf("TestLeakyBucketRateLimiter() Elapsed time should have been around 2000ms, but it was %d", elapsed)
		}
	})
}
