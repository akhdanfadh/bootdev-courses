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

func (c *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get Bearer token from the request headers
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJson(w, http.StatusUnauthorized, errorResponse{Error: "No access token provided"})
		return
	}
	// Validate access token and get user ID
	userID, err := auth.ValidateJWT(token, c.JwtSecret)
	if err != nil {
		respondJson(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}

	// Decode the request
	type validRequest struct { // JSON structs for request
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	request := validRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Crucial to match the validRequest struct
	if err := decoder.Decode(&request); err != nil {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON request"})
		return
	}

	// Validate request
	newEmail := strings.ToLower(strings.TrimSpace(request.Email))
	newPassword := request.Password
	if newEmail == "" && newPassword == "" {
		respondJson(w, http.StatusBadRequest, errorResponse{Error: "Provide either email or password or both"})
		return
	}
	// If updating email, validate it
	if newEmail != "" {
		if _, err := mail.ParseAddress(newEmail); err != nil {
			respondJson(w, http.StatusBadRequest, errorResponse{Error: "Invalid email format"})
			return
		}
	}
	// If updating password, hash it
	if newPassword != "" {
		newPassword, err = auth.HashPassword(newPassword)
		if err != nil {
			respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error: failed to hash password"})
			return
		}
	}

	// Update the database
	user, err := c.db.UpdateUserEmailPassword(r.Context(), database.UpdateUserEmailPasswordParams{
		ID:             userID,
		Email:          newEmail,
		HashedPassword: newPassword,
	})
	if err != nil {
		// Check for duplicate, since email is unique
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			respondJson(w, http.StatusConflict, errorResponse{Error: "User with this email already exists"})
			return
		}
		log.Println("Error changing user in database:", err)
		respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error"})
		return
	}

	// If all good, return the user data
	respondJson(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
