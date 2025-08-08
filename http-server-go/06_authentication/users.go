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

// handlerLogin is an HTTP handler function to handle user login
func (c *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Constants for token expiration
	const accessTokenExpires = time.Duration(1) * time.Hour
	const refreshTokenExpires = time.Duration(24*60) * time.Hour

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

	// Get the user email from database
	user, err := c.db.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		respondJson(w, http.StatusNotFound, errorResponse{Error: "Email not found"})
		return
	}

	// Compare the password
	err = auth.CheckPasswordHash(request.Password, user.HashedPassword)
	if err != nil {
		respondJson(w, http.StatusUnauthorized, errorResponse{Error: "Incorrect email or password"})
		return
	}

	// Create refresh token and store it in the database
	refreshToken, err := c.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(refreshTokenExpires),
	})
	if err != nil {
		respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error: failed to create refresh token"})
		return
	}

	// Create JWT
	accessToken, err := auth.MakeJWT(user.ID, c.JwtSecret, accessTokenExpires)
	if err != nil {
		respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error: failed to create JWT token"})
		return
	}

	// If all good, return the user data
	respondJson(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})
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
