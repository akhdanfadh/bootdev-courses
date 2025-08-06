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

	// A simple fileserver on current directory (./index.html)
	fileServerHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/", fileServerHandler)

	// A simple way to run HTTP server with configured parameters
	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
