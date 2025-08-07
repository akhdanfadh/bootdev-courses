package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// handlerValidateChirp is a function to handle validating chirp post endpoint
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	// JSON structs for request and responses
	type validRequest struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		Body string `json:"cleaned_body"`
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

// cleanChirp is a function to clean the chirp text by replacing banned words with ****
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
