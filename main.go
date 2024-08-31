package main

import (
	"context"
	"fmt"
	"go-socket/internal/db"
	"log"
	"net/http"

	apikey "go-socket/internal/models"

	"github.com/gorilla/websocket"
)

var CurrentConnectedUsers []apikey.APIKey

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("api_key")
	if apiKey == "" {
		http.Error(w, "API Key é obrigatória", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	rdb := db.NewDatabaseClient()

	exists, err := rdb.Exists(ctx, apiKey).Result()
	if err != nil {
		http.Error(w, "Erro ao verificar API Key", http.StatusInternalServerError)
		return
	}

	if exists == 0 {
		http.Error(w, "API Key inválida", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Erro ao ler mensagem:", err)
			break
		}
		fmt.Printf("Mensagem recebida: %s\n", message)

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Erro ao enviar mensagem:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Servidor WebSocket rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
