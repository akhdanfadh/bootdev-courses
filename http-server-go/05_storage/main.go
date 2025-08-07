package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// Load config from environment variables
	cfg := loadConfig()

	// Open connection to database
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// Stateful configuration for the API
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
	}

	// A multiplexer is responsible for routing HTTP requests to appropriate handler
	mux := http.NewServeMux()

	appHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(cfg.FilepathRoot))))
	mux.Handle("/app/", appHandler) // fileserver on current directory as '/app' endpoint

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics) // show the fileserverHits
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)    // reset the fileserverHits

	mux.HandleFunc("GET /api/healthz", handlerReadiness)             // healthcheck endpoint
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp) // validate chirp endpoint

	// A simple way to run HTTP server with configured parameters
	// The use of pointer is to avoid accidental copies when passing between func/goroutines
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + cfg.Port,
	}
	log.Printf("Serving files from %s on port: %s\n", cfg.FilepathRoot, cfg.Port)
	log.Fatal(server.ListenAndServe())
}
