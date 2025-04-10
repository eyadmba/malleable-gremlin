package postgresql

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/eyadmba/malleable-gremlin/services/postgresql"
)

// SetupRoutes configures routes for the PostgreSQL service
func SetupRoutes(prefix string, router *http.ServeMux, pgManager *postgresql.ConnectionManager) {
	router.HandleFunc("PUT "+prefix+"/connection-string", handleStoreConnectionString(pgManager))
	router.HandleFunc("POST "+prefix+"/connect", handlePostgresConnect(pgManager))
	router.HandleFunc("POST "+prefix+"/query", handleExecuteQuery(pgManager))
}

// handleStoreConnectionString returns a handler for the PostgreSQL connection string endpoint
func handleStoreConnectionString(pgManager *postgresql.ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ConnectionString string `json:"connectionString"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		result := pgManager.StoreConnectionString(req.ConnectionString)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// handlePostgresConnect returns a handler for the PostgreSQL connect endpoint
func handlePostgresConnect(pgManager *postgresql.ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ConnectionString   string `json:"connectionString"`
			ConnectionStringID string `json:"connectionStringId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		response, err := pgManager.Connect(r.Context(), req.ConnectionString, req.ConnectionStringID)
		if err != nil {
			if errors.Is(err, postgresql.ErrBothInputsProvided) ||
				errors.Is(err, postgresql.ErrNeitherInputProvided) ||
				errors.Is(err, postgresql.ErrConnIDNotFound) ||
				errors.Is(err, postgresql.ErrConnectionSetupFailed) {
				http.Error(w, err.Error(), http.StatusBadRequest) // 400 Bad Request
			} else if errors.Is(err, postgresql.ErrConnectionFailed) {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			} else {
				log.Printf("Internal server error during connect: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// handleExecuteQuery returns a handler for the PostgreSQL query endpoint
func handleExecuteQuery(pgManager *postgresql.ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req postgresql.ExecuteQueryArgument
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		result, err := pgManager.ExecuteQuery(r.Context(), &req)
		if err != nil {
			if errors.Is(err, postgresql.ErrBothInputsProvided) ||
				errors.Is(err, postgresql.ErrNeitherInputProvided) ||
				errors.Is(err, postgresql.ErrConnIDNotFound) ||
				errors.Is(err, postgresql.ErrConnectionSetupFailed) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else if errors.Is(err, postgresql.ErrConnectionFailed) {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			} else {
				log.Printf("Internal server error during query execution setup: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
