package redis

import (
	"context"
	
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisLocker struct {
	client *redis.Client
}

func NewRedisLocker(client *redis.Client) *RedisLocker {
	return &RedisLocker{client: client}
}
func (r *RedisLocker) LockResource(
	ctx context.Context,
	resource string,
	id uint,
	expiry time.Duration,
) (func(), error) {

	key := fmt.Sprintf("%s:%d", resource, id)
	token := uuid.NewString()

	res, err := r.client.SetArgs(ctx, key, token, redis.SetArgs{
		Mode: "NX",
		TTL:  expiry,
	}).Result()
	if err != nil {
		return nil, err
	}

	if res != "OK" {
		return nil, fmt.Errorf("resource %s:%d is locked", resource, id)
	}

	//  heartbeat control
	stopChan := make(chan struct{})

	//  Start TTL refresh goroutine
	go func() {
		ticker := time.NewTicker(expiry / 2) // refresh before expiry
		defer ticker.Stop()

		script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end
		`

		for {
			select {
			case <-ticker.C:
				_, err := r.client.Eval(
					context.Background(),
					script,
					[]string{key},
					token,
					int(expiry.Milliseconds()),
				).Result()

				if err != nil {
					// optional: log error
				}

			case <-stopChan:
				return
			}
		}
	}()
		unlock := func() {
		close(stopChan) // stop heartbeat

		script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
		`

		_, err := r.client.Eval(
			context.Background(),
			script,
			[]string{key},
			token,
		).Result()

		if err != nil {
			fmt.Printf("failed to release lock: %v\n", err)
		}
	}

	return unlock, nil
}