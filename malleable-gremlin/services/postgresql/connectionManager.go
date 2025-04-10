package postgresql

import (
	"sync"
)

type ConnectionManager struct {
	connections map[string]string
	mu          sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]string),
	}
}

func (cm *ConnectionManager) Close() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.connections = make(map[string]string)
}
