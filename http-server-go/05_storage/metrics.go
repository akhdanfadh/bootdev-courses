package main

import (
	"fmt"
	"net/http"
)

// middlewareMetricsInc is a middleware function that increments the fileserverHits counter
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// http.Handler is an interface that has ServeHTTP method
	// http.HandlerFunc is a type conversion that gives the function that method
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r) // execute the original handler
	})
}

// handlerMetrics is a function to handle showing the metrics of the fileserverHits
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
