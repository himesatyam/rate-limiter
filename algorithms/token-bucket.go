package algorithms

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rate-limiter/config"
	internal_redis "github.com/rate-limiter/redis"
	"github.com/redis/go-redis/v9"
)

var memstore map[string]store

type tokenBucketLimiter struct {
	config config.Config
	redis  *redis.Client
}

type store struct {
	size       int
	timeStored time.Time
}

type redisStore struct {
	size int `redis:"size"`
}

func NewTokenBucketLimiter(config config.Config) Algorithms {
	limiter := tokenBucketLimiter{}
	if config.Redis.Enable {
		limiter.redis = internal_redis.NewRedisClient(config.Redis.Options)
	}
	limiter.config = config
	memstore = make(map[string]store)

	if !config.Redis.Enable {
		go limiter.clearStore()
	}
	return limiter
}

func (l tokenBucketLimiter) Allow(key string) (bool, error) {
	if l.config.Redis.Enable {
		return l.allowWithRedis(key)
	}
	val, ok := memstore[key]
	if !ok {
		memstore[key] = store{
			size:       l.config.TokenBucket.BucketSize - 1,
			timeStored: time.Now(),
		}
		return true, nil
	}

	if val.size > 0 && time.Since(val.timeStored) < l.config.TokenBucket.UpdateEvery {
		val.size--
		memstore[key] = val
		return true, nil
	}

	return false, nil
}

func (l *tokenBucketLimiter) allowWithRedis(key string) (bool, error) {
	var s redisStore
	result, err := l.redis.HGetAll(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	if len(result) == 0 {
		s.size = l.config.TokenBucket.BucketSize - 1
		err = l.redis.HSet(context.Background(), key, "size", s.size).Err()
		if err != nil {
			return false, err
		}
		err = l.redis.Expire(context.Background(), key, l.config.TokenBucket.UpdateEvery).Err()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	size, err := strconv.Atoi(result["size"])
	if err != nil {
		return false, err
	}

	if size <= 0 {
		return false, nil
	}

	size--
	err = l.redis.HSet(context.Background(), key, "size", size).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (l tokenBucketLimiter) clearStore() {
	for {
		time.Sleep(l.config.TokenBucket.UpdateEvery)
		for k, s := range memstore {
			if time.Since(s.timeStored).Seconds() >= l.config.TokenBucket.UpdateEvery.Seconds() {
				delete(memstore, k)
			}
		}
		fmt.Println("Store cleared")
	}
}
