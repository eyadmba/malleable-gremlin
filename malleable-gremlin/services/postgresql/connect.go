package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var connectTimeout = 10 * time.Second

type ConnectResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (cm *ConnectionManager) Connect(ctx context.Context, connStr, connID string) (*ConnectResult, error) {
	if connStr == "" && connID == "" {
		return nil, ErrNeitherInputProvided
	}

	if connStr != "" && connID != "" {
		return nil, ErrBothInputsProvided
	}

	var db *sql.DB
	var err error
	var actualConnStr string

	if connStr != "" {
		actualConnStr = connStr
	} else {
		cm.mu.RLock()
		storedConnStr, exists := cm.connections[connID]
		cm.mu.RUnlock()

		if !exists {
			return nil, fmt.Errorf("%w: '%s'", ErrConnIDNotFound, connID)
		}
		actualConnStr = storedConnStr
	}

	// Open a new connection using the determined string
	db, err = sql.Open("postgres", actualConnStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectionSetupFailed, err)
	}
	defer db.Close()

	// Test the connection using a derived context with timeout
	pingCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return &ConnectResult{
			Success: false,
			Error:   fmt.Errorf("%w: %w", ErrConnectionFailed, err).Error(),
		}, nil
	}

	return &ConnectResult{
		Success: true,
	}, nil
}
