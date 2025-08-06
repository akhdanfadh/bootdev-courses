package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// Struct to hold any stateful, in-memory data across requests
type apiConfig struct {
	fileserverHits atomic.Int32 // atomic allows to safely use value across goroutines
}

// Middleware that increments fileserverHits every time it's called
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// http.Handler is an interface that has ServeHTTP method
	// http.HandlerFunc is a type conversion that gives the function that method
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r) // execute the original handler
	})
}

func main() {
	// Server configuration
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// A multiplexer is responsible for routing HTTP requests to appropriate handler
	mux := http.NewServeMux()
	// A simple fileserver on current directory (./index.html) on /app endpoint
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	// Endpoint for showing the fileserverHits
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	// Endpoint to reset the fileserverHits
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)
	// A custom handler for readiness endpoint
	mux.HandleFunc("GET /healthz", handlerReadiness)

	// A simple way to run HTTP server with configured parameters
	// The use of pointer is to avoid accidental copies when passing between func/goroutines
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

// A custom function to handle the readiness endpoint, simply return 200 OK
func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK))) // just return 200 OK
}

// A custom function to handle the hits endpoint
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())
}

// A custom function to handle resetting the fileserverHits
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
