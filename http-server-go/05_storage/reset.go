package main

import (
	"fmt"
	"log"
	"net/http"
)

// handlerReset handles the reset endpoint.
// It resets the fileserverHits and deletes all user data from the database.
func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// Only allow on "dev" platform
	if c.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden: Reset is only allowed in development environment."))
		return
	}

	// Reset fileserverHits
	c.fileserverHits.Store(0)

	// Reset users data
	err := c.db.DeleteAllUsers(r.Context())
	if err != nil {
		text := fmt.Sprintf("Error deleting database: %s", err)
		log.Println(text)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(text))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database and metrics have been reset."))
}
