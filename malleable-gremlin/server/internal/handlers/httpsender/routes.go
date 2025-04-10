package httpsender

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/eyadmba/malleable-gremlin/services/httpsender"
)

func SetupRoutes(prefix string, router *http.ServeMux) {
	router.HandleFunc("GET "+prefix+"/send/{forwardUrl...}", handleForwardWhenGet)
	router.HandleFunc("POST "+prefix+"/send", handleHTTPSend)
}

func handleForwardWhenGet(w http.ResponseWriter, r *http.Request) {
	forwardUrl := r.PathValue("forwardUrl")
	domain, path, _ := strings.Cut(forwardUrl, "/")

	u := &url.URL{
		Scheme:   "http",
		Host:     domain,
		Path:     path,
		RawQuery: r.URL.RawQuery,
	}

	req := &httpsender.SendArgument{
		URL:    u.String(),
		Method: "GET",
	}

	response, err := httpsender.Send(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func handleHTTPSend(w http.ResponseWriter, r *http.Request) {
	var req httpsender.SendArgument
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := httpsender.Send(&req)
	if err != nil {
		log.Printf("HTTP send failed: %v", err)
		http.Error(w, "failed to send HTTP request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.StatusCode)
	json.NewEncoder(w).Encode(result.Body)
}
