package apikey

import "time"

type KeyMetaData struct {
	Name          string
	CreatedAt     time.Time
	LastConnected time.Time
	IsConnected   bool
}

type APIKey struct {
	Key  string
	Data KeyMetaData
}
