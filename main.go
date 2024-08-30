package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// var redisClient = redis.NewClient(&redis.Options{
// 	Addr:     "localhost:6379",
// 	Password: "",
// 	DB:       1,
// })

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

func ProcessClient() {

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       1,
	})

	err := rdb.HSet(ctx, "user:1000", "name", "Luck Ramos", "age", "22").Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.HGet(ctx, "user:1000", "name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("user:1000", val)
}

// func main() {
// 	ProcessClient()
// }

func main() {
	http.HandleFunc("/ws", handlers.HandleConnections)
	http.HandleFunc("/login", handlers.Login)

	// Inicia o servidor
	log.Println("Servidor rodando na porta 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
