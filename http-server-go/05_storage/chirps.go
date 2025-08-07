package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

// bannedWordsMap is a set of banned words (Go does not have Set as in Python)
var bannedWordsMap = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

// handlerAddChirp is an HTTP handler function to add a chirp
func (c *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	// JSON structs for request and responses
	type validRequest struct {
		UserID uuid.UUID `json:"user_id"`
		Body   string    `json:"body"`
	}

	// Decode the request
	request := validRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Crucial to match the validRequest struct
	if err := decoder.Decode(&request); err != nil {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON request"})
		return
	}

	// Validate and clean chirp
	if request.UserID == uuid.Nil || request.Body == "" {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "User ID and Body are required"})
		return
	}
	const maxChirpLength = 140
	if len(request.Body) > maxChirpLength {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Chirp is too long"})
		return
	}
	cleaned := cleanChirp(request.Body)

	// Store on database
	chirp, err := c.db.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: request.UserID,
		Body:   cleaned,
	})
	if err != nil {
		// Check for invalid user ID
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			respondJson(w, http.StatusBadRequest, errorResponse{Error: "Invalid user ID"})
			return
		}
		log.Println("Error adding chirp to database:", err)
		respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error"})
		return
	}

	// If all good, return the chirp data
	respondJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
	})
}

// cleanChirp is a function to clean the chirp text by replacing banned words with ****
func cleanChirp(chirp string) string {
	// Split on whitespace, change banned to ****, then join
	words := strings.Split(chirp, " ")
	for i, word := range words {
		if bannedWordsMap[strings.ToLower(word)] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
