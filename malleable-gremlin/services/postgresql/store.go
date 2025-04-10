package postgresql

import (
	"fmt"
	"time"
)

type StoreConnectionStringResult struct {
	ID string `json:"id"`
}

func (cm *ConnectionManager) StoreConnectionString(connStr string) *StoreConnectionStringResult {
	// Generate a unique ID for this connection string
	id := fmt.Sprintf("conn_%d", time.Now().UnixNano())

	cm.mu.Lock()
	cm.connections[id] = connStr
	cm.mu.Unlock()

	return &StoreConnectionStringResult{ID: id}
}
