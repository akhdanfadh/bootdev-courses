package main

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/database"
	"github.com/joho/godotenv"
)

// envConfig holds the application configuration
type envConfig struct {
	DBUrl, Port, FilepathRoot, Platform, JwtSecret, PolkaKey string
}

// apiConfig holds the stateful configuration for the API
type apiConfig struct {
	fileserverHits atomic.Int32      // atomic allows to safely use value across goroutines
	db             *database.Queries // sqlc-generated-struct to interact with the database
	platform       string            // environment of running application, e.g. "dev", "prod"
	JwtSecret      string            // secret key for JWT signing
	PolkaKey       string            // key for Polka webhook
}

// validate checks if all required applcication configuration fields are set
func (e *envConfig) validate() {
	envVars := map[string]string{
		"DB_URL":           e.DBUrl,
		"CHIRPY_PORT":      e.Port,
		"CHIRPY_FILE_ROOT": e.FilepathRoot,
		"PLATFORM":         e.Platform,
		"JWT_SECRET":       e.JwtSecret,
		"POLKA_KEY":        e.PolkaKey,
	}

	for envName, envValue := range envVars {
		if envValue == "" {
			log.Fatalf("%s must be set in .env", envName)
		}
	}
}

// loadEnv loads and validates application configuration from environment variables
func loadEnv() *envConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := &envConfig{
		DBUrl:        os.Getenv("DB_URL"),
		Port:         os.Getenv("CHIRPY_PORT"),
		FilepathRoot: os.Getenv("CHIRPY_FILE_ROOT"),
		Platform:     os.Getenv("PLATFORM"),
		JwtSecret:    os.Getenv("JWT_SECRET"),
		PolkaKey:     os.Getenv("POLKA_KEY"),
	}
	env.validate()

	return env
}
