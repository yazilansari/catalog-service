package redis

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *redis.Client

func ConnectRedis() {
	db, _ := strconv.Atoi(
		os.Getenv("REDIS_DB"),
	)

	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
}
