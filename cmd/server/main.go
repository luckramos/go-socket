package main

import (
	"context"
	"fmt"
	"go-socket/internal/db"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var currentConnectedClients = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
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

	currentConnectedClients[apiKey] = conn

	// ticker := time.NewTicker(2 * time.Second)
	// quit := make(chan struct{})
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			teste := "ping"

	// 			sendMessageToClient(apiKey, []byte(teste))

	// 		case <-quit:
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(currentConnectedClients, apiKey)
			conn.Close()
			break
		}

		sendMessageToClient(apiKey, msg)
	}

}

func sendMessageToClient(apiKey string, message []byte) {
	if ws, ok := currentConnectedClients[apiKey]; ok {
		ws.WriteMessage(websocket.TextMessage, message)
	}
}

func main() {
	fmt.Println("Servidor WebSocket rodando na porta 8080...")
	http.HandleFunc("/", handleConnections)
	http.ListenAndServe(":8080", nil)
}
