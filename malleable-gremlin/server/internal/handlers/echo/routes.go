package echo

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// SetupRoutes configures routes for the Echo service
func SetupRoutes(prefix string, router *http.ServeMux) {
	router.HandleFunc("GET "+prefix+"/get", handleGetEcho)
	router.HandleFunc("POST "+prefix+"/post", handlePostEcho)
}

// handleGetEcho echoes back request details for GET requests
func handleGetEcho(w http.ResponseWriter, r *http.Request) {
	status := 200
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		if s, err := strconv.Atoi(statusParam); err == nil {
			status = s
		}
	}

	response := map[string]interface{}{
		"args":    r.URL.Query(),
		"headers": r.Header,
		"url":     r.URL.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// handlePostEcho echoes back request details for POST requests
func handlePostEcho(w http.ResponseWriter, r *http.Request) {
	status := 200
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		if s, err := strconv.Atoi(statusParam); err == nil {
			status = s
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"args":    r.URL.Query(),
		"headers": r.Header,
		"url":     r.URL.String(),
		"form":    r.PostForm,
		"files":   r.MultipartForm,
		"data":    string(body),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
