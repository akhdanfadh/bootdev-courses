package main

import (
	"log"
	"net/http"
)

func main() {
	// Server configuration
	const port = "8080"
	const filepathRoot = "."

	// A multiplexer is responsible for routing HTTP requests to appropriate handler
	mux := http.NewServeMux()

	// A simple fileserver on current directory (./index.html) on /app endpoint
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	// A custom handler for readiness endpoint
	mux.HandleFunc("/healthz", handlerReadiness)

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
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)                    // just return 200 OK
	w.Write([]byte(http.StatusText(http.StatusOK))) // response body
}
