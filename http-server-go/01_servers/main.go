package main

import "net/http"

func main() {
	// A multiplexer is responsible for routing HTTP requests to appropriate handler
	mux := http.NewServeMux()

	// A simple fileserver on current directory (./index.html)
	fileServerHandler := http.FileServer(http.Dir("."))
	// a url with pattern /app/ will directed to the fileserver (which is cwd)
	mux.Handle("/", fileServerHandler)

	// A simple way to run HTTP server with configured parameters
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
