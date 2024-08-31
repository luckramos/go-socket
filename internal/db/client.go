package db

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func NewDatabaseClient() *redis.Client {

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	return rdb
}
