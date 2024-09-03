package apikey

import (
	"time"
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
