package load

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/eyadmba/malleable-gremlin/services/load"
)

// SetupRoutes configures routes for the Load service
func SetupRoutes(prefix string, router *http.ServeMux) {
	router.HandleFunc("GET "+prefix+"/cpu", handleCPULoad())
	router.HandleFunc("GET "+prefix+"/memory", handleMemoryLoad())
	router.HandleFunc("GET "+prefix+"/io", handleIOLoad())
}

type LoadResponse struct {
	TasksStarted int    `json:"tasks_started"`
	Duration     string `json:"duration"`
	Error        string `json:"error"`
}

// handleCPULoad returns a handler for the CPU load endpoint
func handleCPULoad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasksStr := r.URL.Query().Get("tasks")
		timeoutStr := r.URL.Query().Get("timeout")

		if timeoutStr == "" {
			http.Error(w, "timeout parameter is required", http.StatusBadRequest)
			return
		}

		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			http.Error(w, "invalid timeout format", http.StatusBadRequest)
			return
		}

		var tasks int
		if tasksStr == "cpus" {
			tasks = runtime.NumCPU()
		} else {
			tasks, err = strconv.Atoi(tasksStr)
			if err != nil {
				http.Error(w, "invalid tasks format", http.StatusBadRequest)
				return
			}
		}

		start := time.Now()
		result, err := load.GenerateCPULoad(r.Context(), tasks, timeout)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		elapsed := time.Since(start)

		response := &LoadResponse{
			Error:        result.Error,
			TasksStarted: result.TasksStarted,
			Duration:     fmt.Sprintf("%dms", elapsed.Milliseconds()),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// handleMemoryLoad returns a handler for the memory load endpoint
func handleMemoryLoad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sizeStr := r.URL.Query().Get("size")
		gcAfterStr := r.URL.Query().Get("gc_after")

		if sizeStr == "" {
			http.Error(w, "size parameter is required", http.StatusBadRequest)
			return
		}

		// Parse size (e.g., "500mb", "1.5gb")
		size, err := parseSize(sizeStr)
		if err != nil {
			http.Error(w, "invalid size format", http.StatusBadRequest)
			return
		}

		var gcAfter time.Duration
		if gcAfterStr != "" {
			if gcAfterStr == "-1" {
				gcAfter = -1
			} else {
				gcAfter, err = time.ParseDuration(gcAfterStr)
				if err != nil {
					http.Error(w, "invalid gc_after format", http.StatusBadRequest)
					return
				}
			}
		}

		result, err := load.GenerateMemoryLoad(size, gcAfter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// handleIOLoad returns a handler for the I/O load endpoint
func handleIOLoad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasksStr := r.URL.Query().Get("tasks")
		waitStr := r.URL.Query().Get("wait")
		parallelStr := r.URL.Query().Get("parallel")

		if tasksStr == "" || waitStr == "" || parallelStr == "" {
			http.Error(w, "tasks, wait, and parallel parameters are required", http.StatusBadRequest)
			return
		}

		tasks, err := strconv.Atoi(tasksStr)
		if err != nil {
			http.Error(w, "invalid tasks format", http.StatusBadRequest)
			return
		}

		wait, err := time.ParseDuration(waitStr)
		if err != nil {
			http.Error(w, "invalid wait format", http.StatusBadRequest)
			return
		}

		parallel, err := strconv.Atoi(parallelStr)
		if err != nil {
			http.Error(w, "invalid parallel format", http.StatusBadRequest)
			return
		}

		// Pass the request context to the service function
		result, err := load.GenerateIOLoad(r.Context(), tasks, wait, parallel)
		if err != nil {
			// Check if the error is due to context cancellation
			if err == context.Canceled || err == context.DeadlineExceeded {
				// Return a specific status code for cancellation, e.g., 499 Client Closed Request
				// or handle it as a server error, depending on requirements.
				http.Error(w, fmt.Sprintf("Request cancelled: %v", err), 499)
			} else {
				// Handle other errors as internal server errors
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// parseSize parses a size string (e.g., "500mb", "1.5gb") into bytes
func parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToLower(strings.TrimSpace(sizeStr))
	var multiplier int64 = 1

	if strings.HasSuffix(sizeStr, "kb") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "kb")
	} else if strings.HasSuffix(sizeStr, "mb") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "mb")
	} else if strings.HasSuffix(sizeStr, "gb") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "gb")
	} else {
		return 0, fmt.Errorf("invalid size unit")
	}

	sizeStr = strings.TrimSpace(sizeStr)
	sizeFloat, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, err
	}

	return int64(sizeFloat * float64(multiplier)), nil
}
