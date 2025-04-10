package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/eyadmba/malleable-gremlin/server/internal"
	"github.com/eyadmba/malleable-gremlin/server/internal/handlers/about"
	"github.com/eyadmba/malleable-gremlin/server/internal/handlers/echo"
	"github.com/eyadmba/malleable-gremlin/server/internal/handlers/httpsender"
	"github.com/eyadmba/malleable-gremlin/server/internal/handlers/load"
	"github.com/eyadmba/malleable-gremlin/server/internal/handlers/postgresql"
)

// readArgs parses command line arguments
func readArgs() string {
	addr := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()
	return *addr
}

// setupServer creates and configures the HTTP server
func setupServer(deps *internal.Dependencies) *http.ServeMux {
	router := http.NewServeMux()

	// Setup routes for each service with its prefix
	about.SetupRoutes("/about", router)
	echo.SetupRoutes("/echo", router)
	load.SetupRoutes("/load", router)
	httpsender.SetupRoutes("/http-send", router)
	postgresql.SetupRoutes("/postgresql", router, deps.PostgresManager)

	return router
}

// runServer starts the HTTP server and handles shutdown
func runServer(addr string, router *http.ServeMux, deps *internal.Dependencies) {
	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to handle server errors
	errChan := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on %s", addr)
		errChan <- http.ListenAndServe(addr, router)
	}()

	// Wait for either a signal or an error
	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
		os.Exit(1)
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down...", sig)
		deps.Close()
	}
}

func main() {
	// Parse command line arguments
	addr := readArgs()

	// Build dependencies
	deps := internal.BuildDependencies()

	// Setup and configure the server
	router := setupServer(deps)

	// Run the server
	runServer(addr, router, deps)
}
