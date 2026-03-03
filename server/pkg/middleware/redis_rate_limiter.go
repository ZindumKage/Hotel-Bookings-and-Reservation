package middleware

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)



type RedisRateLimiter struct {
	client *redis.Client
}

func (r *RedisRateLimiter) Allow(key string, limit int, window time.Duration) (bool, error) {
	ctx := context.Background()

	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.client.Expire(ctx, key, window)
	}

	return count <= int64(limit), nil
}