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
	env := loadEnv()

	// Open connection to database
	db, err := sql.Open("postgres", env.DBUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// Stateful configuration for the API
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       env.Platform,
		JwtSecret:      env.JwtSecret,
	}

	// A multiplexer is responsible for routing HTTP requests to appropriate handler
	mux := http.NewServeMux()

	appHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(env.FilepathRoot))))
	mux.Handle("/app/", appHandler) // fileserver on current directory as '/app' endpoint

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics) // show the fileserverHits
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)    // reset the fileserverHits and database

	mux.HandleFunc("GET /api/healthz", handlerReadiness) // healthcheck endpoint

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)          // get all chirps
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp) // get a chirps
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerAddChirp)          // validate chirp endpoint

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)   // login endpoint
	mux.HandleFunc("POST /api/users", apiCfg.handlerAddUser) // add users by email

	// A simple way to run HTTP server with configured parameters
	// The use of pointer is to avoid accidental copies when passing between func/goroutines
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + env.Port,
	}
	log.Printf("Serving files from %s on port: %s\n", env.FilepathRoot, env.Port)
	log.Fatal(server.ListenAndServe())
}
