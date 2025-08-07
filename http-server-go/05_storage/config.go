package main

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/database"
	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DBUrl        string
	Port         string
	FilepathRoot string
}

// apiConfig holds the stateful configuration for the API
type apiConfig struct {
	fileserverHits atomic.Int32      // atomic allows to safely use value across goroutines
	db             *database.Queries // sqlc-generated-struct to interact with the database
}

// validate checks if all required applcication configuration fields are set
func (c *Config) validate() {
	envVars := map[string]string{
		"DB_URL":           c.DBUrl,
		"CHIRPY_PORT":      c.Port,
		"CHIRPY_FILE_ROOT": c.FilepathRoot,
	}

	for envName, envValue := range envVars {
		if envValue == "" {
			log.Fatalf("%s must be set", envName)
		}
	}
}

// loadConfig loads and validates application configuration from environment variables
func loadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{
		DBUrl:        os.Getenv("DB_URL"),
		Port:         os.Getenv("CHIRPY_PORT"),
		FilepathRoot: os.Getenv("CHIRPY_FILE_ROOT"),
	}
	config.validate()

	return config
}
