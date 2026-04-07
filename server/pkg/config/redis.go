package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)


var Ctx = context.Background()
var Redis *redis.Client

func ConnectRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})
	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	log.Println("Redis Connected Successfully")
}

