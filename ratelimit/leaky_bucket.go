package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/df424/gotilities/math"
)

type LeakyBucketRateLimiter struct {
	RequestsPerSecond float32
	requestsRemaining int
	maxRequests       int
	Canceled          bool
	mutex             sync.Mutex
	pendingRequests   []chan bool
}

// NewLeakyBucketRateLimiter ... Construct a new LeakyBucketRateLimiter.
func NewLeakyBucketRateLimiter(ctx context.Context, reqsPerSecond float32, maxRequests int) *LeakyBucketRateLimiter {
	rl := LeakyBucketRateLimiter{
		RequestsPerSecond: reqsPerSecond,
		requestsRemaining: 0,
		maxRequests:       maxRequests,
		Canceled:          false,
		mutex:             sync.Mutex{},
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(1000/rl.RequestsPerSecond))

		for {
			select {
			case <-ticker.C:
				rl.updateBucket()
			case <-ctx.Done():
				rl.mutex.Lock()
				rl.Canceled = true
				rl.mutex.Unlock()

				for _, pendingReq := range rl.pendingRequests {
					pendingReq <- false
				}
				return
			}
		}
	}()

	return &rl
}

func (rl *LeakyBucketRateLimiter) updateBucket() {

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.Canceled {
		panic("Caller tried to get a token from a cancelled rate limiter")
	}

	rl.requestsRemaining = math.MinInt(rl.requestsRemaining+1, rl.maxRequests)

	for rl.requestsRemaining > 0 && len(rl.pendingRequests) > 0 {
		pendingReq := rl.pendingRequests[0]
		rl.pendingRequests = rl.pendingRequests[1:]
		rl.requestsRemaining -= 1
		pendingReq <- true
	}
}

func (rl *LeakyBucketRateLimiter) WaitForToken() bool {
	done := make(chan bool, 1)
	rl.mutex.Lock()

	if rl.Canceled {
		rl.mutex.Unlock()
		panic("Caller tried to get a token from a cancelled rate limiter")
	}

	if rl.requestsRemaining > 0 {
		done <- true
		rl.requestsRemaining -= 1
	} else {
		rl.pendingRequests = append(rl.pendingRequests, done)
	}

	rl.mutex.Unlock()
	return <-done
}
