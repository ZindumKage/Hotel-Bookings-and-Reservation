package user

import "time"


type RedisRateLimiter interface {
	Allow(key string, limit int, window time.Duration) (bool, error)
}