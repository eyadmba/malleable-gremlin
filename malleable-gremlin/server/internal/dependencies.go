package internal

import (
	"github.com/eyadmba/malleable-gremlin/services/postgresql"
)

// Dependencies holds all the dependencies needed by the server
type Dependencies struct {
	PostgresManager *postgresql.ConnectionManager
	// Add more dependencies here as needed
}

// BuildDependencies creates and initializes all required dependencies
func BuildDependencies() *Dependencies {
	// Create PostgreSQL connection manager
	pgManager := postgresql.NewConnectionManager()

	// Return all dependencies
	return &Dependencies{
		PostgresManager: pgManager,
		// Add more dependencies here as needed
	}
}

// Close closes all resources held by dependencies
func (d *Dependencies) Close() {
	if d.PostgresManager != nil {
		d.PostgresManager.Close()
	}
	// Close other dependencies here as needed
} 