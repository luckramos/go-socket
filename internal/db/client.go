package db

import (
	"context"
	"encoding/json"
	"fmt"
	apikey "go-socket/internal/models"
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

func RetrieveApiKey(ctx context.Context, rdb *redis.Client, apiKey string) (*apikey.APIKey, error) {

	apiKeyJSON, err := rdb.HGet(ctx, apiKey, "data").Result()
	if err != nil {
		return nil, err
	}

	var apiUser apikey.APIKey
	err = json.Unmarshal([]byte(apiKeyJSON), &apiUser)
	if err != nil {
		return nil, fmt.Errorf("erro ao desserializar API Key: %v", err)
	}

	return &apiUser, nil
}
