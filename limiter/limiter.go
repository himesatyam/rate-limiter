package limiter

import (
	"github.com/rate-limiter/algorithms"
	"github.com/rate-limiter/config"
)

// NewRateLimiter returns rate limiter based on the configuration provided.
// If TokenBucket is enabled in the configuration, it will return TokenBucketLimiter.
// If LeakingBucket is enabled in the configuration, it will return LeakingBucketLimiter.
// If none of the above is enabled, it will return nil.
// If Redis is enabled in the configuration, it will return Redis based rate limiter.
// If Redis is not enabled, it will return in-memory based rate limiter.
func NewRateLimiter(config config.Config) algorithms.Algorithms {
	if config.LimiterAlgorithm.TokenBucket {
		return algorithms.NewTokenBucketLimiter(config)
	} else if config.LimiterAlgorithm.LeakingBucket {
	}

	return nil
}
