package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	// Endpoint to reset the fileserverHits
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	// A custom handler for readiness endpoint
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// Endpoint to validate chirp
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

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
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	htmlTemplate := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	fmt.Fprintf(w, htmlTemplate, cfg.fileserverHits.Load())
}

// A custom function to handle resetting the fileserverHits
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

// handlerValidateChirp is a function to handle validating chirp post endpoint
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	// JSON structs for request and responses
	type validRequest struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		Body string `json:"cleaned_body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	const maxChirpLength = 140

	// Decode the request
	request := validRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON request"})
		return
	}

	// Validate the length
	if len(request.Body) > maxChirpLength {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Chirp is too long"})
		return
	}

	// Then if everything okay
	cleaned := cleanChirp(request.Body)
	respondJson(w, http.StatusOK, validResponse{Body: cleaned})
}

// respondJSON is a utility function to respond with JSON data given payload
func respondJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload) // encode the payload
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(500) // since encoding is server-side
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func cleanChirp(chirp string) string {
	// Make a set of banned words (Go does not have Set as in Python)
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	banned := make(map[string]bool)
	for _, word := range bannedWords {
		banned[word] = true
	}

	// Split on whitespace, change banned to ****, then join
	words := strings.Split(chirp, " ")
	for i, word := range words {
		if banned[strings.ToLower(word)] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
