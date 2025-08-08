package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/auth"
	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"token,omitempty"`         // show this only on login endpoint
	RefreshToken string    `json:"refresh_token,omitempty"` // show this only on login endpoint
}

func (c *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	// JSON structs for request
	type validRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// Decode the request
	request := validRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Crucial to match the validRequest struct
	if err := decoder.Decode(&request); err != nil {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON request"})
		return
	}

	// Validate request
	request.Email = strings.TrimSpace(request.Email)
	if request.Email == "" || request.Password == "" {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Email and Password is required"})
		return
	}
	request.Email = strings.ToLower(request.Email)
	if _, err := mail.ParseAddress(request.Email); err != nil {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Invalid email format"})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error: failed to hash password"})
		return
	}

	// Store on database
	user, err := c.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		// Check for duplicate, since email is unique
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			respondJson(w, http.StatusConflict, errorResponse{Error: "User with this email already exists"})
			return
		}
		log.Println("Error adding user to database:", err)
		respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error"})
		return
	}

	// If all good, return the user data
	response := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondJson(w, http.StatusCreated, response)
}
