package main

import (
	"context"
	"fmt"
	"go-socket/internal/db"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	currentConnectedClients = sync.Map{}
	upgrader                = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("api_key")
	if key == "" {
		http.Error(w, "API Key é obrigatória", http.StatusUnauthorized)
		return
	}

	rdb := db.NewDatabaseClient()
	ctx := context.Background()

	apiUser, err := db.RetrieveApiKey(ctx, rdb, key)
	if err != nil {
		http.Error(w, "Erro ao verificar API Key", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return
	}

	currentConnectedClients.Store(key, conn)
	err = apiUser.UpdateStatus(ctx, rdb, true)
	if err != nil {
		log.Printf("Erro ao atualizar status da API Key: %v", err)
		return
	}

	defer func() {
		err = apiUser.UpdateStatus(ctx, rdb, false)
		if err != nil {
			log.Printf("Erro ao atualizar status da API Key: %v", err)
		}
		currentConnectedClients.Delete(key)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			break
		}

		sendMessageToClient(key, msg)
	}

}

func sendMessageToClient(apiKey string, message []byte) {
	if ws, ok := currentConnectedClients.Load(apiKey); ok {
		conn := ws.(*websocket.Conn)
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func main() {
	fmt.Println("Servidor WebSocket rodando na porta 8080...")
	http.HandleFunc("/", handleConnections)
	http.ListenAndServe(":8080", nil)
}
