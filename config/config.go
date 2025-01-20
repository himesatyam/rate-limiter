package config

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	LimiterAlgorithm struct {
		TokenBucket   bool
		LeakingBucket bool
	}

	Redis struct {
		Enable  bool
		Options redis.Options
	}

	TokenBucket TokenBucket
}

type TokenBucket struct {
	BucketSize  int
	UpdateEvery time.Duration
}

// NewConfig returns config to be used for configuring rate limiter.
// Default rate limiter algorithm will be Token Bucket.
// By default redis is disable if you want to use redis please use Config.Redis and provide required fields
func NewConfig() Config {
	return Config{}
}
