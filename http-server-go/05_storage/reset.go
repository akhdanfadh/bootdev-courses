package main

import (
	"fmt"
	"log"
	"net/http"
)

// handlerReset handles the reset endpoint.
// It resets the fileserverHits and deletes all user data from the database.
func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
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
