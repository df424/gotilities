package ratelimit

type RateLimiter interface {
	WaitForToken() bool
}

type NoLimitRateLimiter struct {
}

func (rl *NoLimitRateLimiter) WaitForToken() bool {
	return true
}
