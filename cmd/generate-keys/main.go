package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"go-socket/internal/db"
	apikey "go-socket/internal/models"
	"log"
	"time"
)

var ctx = context.Background()

func generateAPIKey(name string) *apikey.APIKey {
	key := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}
	encodedKey := hex.EncodeToString(key)

	apiKey := &apikey.APIKey{
		Key: encodedKey,
		Data: apikey.KeyMetaData{
			Name:          name,
			CreatedAt:     time.Now(),
			LastConnected: time.Now(),
			IsConnected:   false,
		},
	}

	return apiKey
}

func main() {
	name := flag.String("name", "API User", "O nome associado à API Key")
	flag.Parse()

	rdb := db.NewDatabaseClient()
	apiKey := generateAPIKey(*name)

	apiKeyJSON, err := json.Marshal(apiKey)
	if err != nil {
		log.Fatalf("Erro ao converter a APIKey para JSON: %v", err)
	}

	err = rdb.HSet(ctx, apiKey.Key, "data", apiKeyJSON).Err()
	if err != nil {
		log.Fatalf("Não foi possível salvar a chave API no Redis: %v", err)
	}

	fmt.Println("API Key gerada e salva com sucesso:", apiKey)
}
