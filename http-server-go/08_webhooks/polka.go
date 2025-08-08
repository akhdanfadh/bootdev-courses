package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/auth"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// Get and validate ApiKey from the Authorization header
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		log.Println("Failed to get Polka API key from request:", err)
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	if apiKey != c.PolkaKey {
		log.Println("Invalid Polka API key in request:", apiKey)
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// Decode the JSON request
	type validRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	request := validRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Crucial to match the validRequest struct
	if err := decoder.Decode(&request); err != nil {
		log.Println("Failed to decode Polka webhook JSON request:", err)
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// Handle the webhook event
	if request.Event != "user.upgraded" {
		log.Println("Unsupported Polka webhook event:", request.Event)
		w.WriteHeader(http.StatusNoContent) // 204
		return
	}

	// Handle the user upgrade event
	userID, err := uuid.Parse(request.Data.UserID)
	if err != nil {
		log.Println("Invalid user ID in Polka webhook request:", request.Data.UserID)
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	_, err = c.db.UpgradeUserChirpyRedByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found in database for upgrade:", userID)
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		log.Println("Failed to upgrade user in database:", err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	// If everything is successful, respond with 204
	log.Println("Successfully upgraded user in database:", userID)
	w.WriteHeader(http.StatusNoContent) // 204
}
