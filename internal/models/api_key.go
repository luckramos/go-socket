package apikey

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type KeyMetaData struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type KeyStatus struct {
	LastConnected time.Time `json:"last_connected"`
	IsConnected   bool      `json:"is_connected"`
}

type APIKey struct {
	Key    string      `json:"key"`
	Meta   KeyMetaData `json:"meta"`
	Status KeyStatus   `json:"status"`
}

func (k *APIKey) UpdateStatus(ctx context.Context, rdb *redis.Client, isConnected bool) error {
	k.Status.IsConnected = isConnected
	k.Status.LastConnected = time.Now()

	apiKeyJSON, err := json.Marshal(k)
	if err != nil {
		return fmt.Errorf("erro ao serializar APIKey: %v", err)
	}

	err = rdb.HSet(ctx, k.Key, "data", apiKeyJSON).Err()
	if err != nil {
		return fmt.Errorf("erro ao atualizar APIKey no banco de dados: %v", err)
	}

	return nil
}
