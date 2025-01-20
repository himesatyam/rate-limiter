package redis

import "github.com/redis/go-redis/v9"

type IRedisClient interface {
}

var client *redis.Client

type RedisClient struct {
}

func NewRedisClient(config redis.Options) *redis.Client {
	if client != nil {
		return client
	}
	return redis.NewClient(&config)
}
