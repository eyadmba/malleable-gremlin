package about

import (
	"encoding/json"
	"net/http"

	"github.com/eyadmba/malleable-gremlin/services/about"
)

// SetupRoutes configures routes for the About service
func SetupRoutes(prefix string, router *http.ServeMux) {
	router.HandleFunc("GET "+prefix+"/system", handleSystemInfo())
	router.HandleFunc("GET "+prefix+"/network", handleNetworkInfo())
}

// handleSystemInfo returns a handler for the system info endpoint
func handleSystemInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := about.GetSystemInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

// handleNetworkInfo returns a handler for the network info endpoint
func handleNetworkInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := about.GetNetworkInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}
