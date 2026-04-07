package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisBlacklist struct {
	client *redis.Client
}

func NewRedisBlacklist(client *redis.Client) *RedisBlacklist {
	return &RedisBlacklist{client: client}
}

func (r *RedisBlacklist) Blacklist(token string, duration time.Duration) error {
	return r.client.Set(context.Background(), token, "revoked", duration).Err()
}

func (r *RedisBlacklist) IsBlacklisted(token string) bool {

	
	val, err := r.client.Get(context.Background(), token).Result()

	if err != nil {
		return false
	}
	return  val == "blacklisted"
}