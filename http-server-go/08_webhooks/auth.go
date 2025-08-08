package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/auth"
	"github.com/akhdanfadh/bootdev-courses/http-server-go/internal/database"
)

// Constants for token expiration
const (
	accessTokenExpires  = time.Duration(1) * time.Hour
	refreshTokenExpires = time.Duration(24*60) * time.Hour
)

// handlerLogin is an HTTP handler function to handle user login
func (c *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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
		ChirpyRed:    user.IsChirpyRed,
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})
}

// handlerRefresh is an HTTP handler function to refresh the access token
// The endpoint does not accept a request body, but expects the refresh token to be sent in the Authorization header
func (c *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// JSON structs for response
	type validResponse struct {
		Token string `json:"token"`
	}

	// Get Bearer bearer from the request headers
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJson(w, http.StatusUnauthorized, errorResponse{Error: "No authorization token provided"})
		return
	}

	// Check if refresh token is valid
	refreshToken, err := c.db.GetRefreshTokenByToken(r.Context(), bearer)
	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		respondJson(w, http.StatusUnauthorized, errorResponse{Error: "Invalid or expired refresh token"})
		return
	}

	// Create new access token (JWT)
	newToken, err := auth.MakeJWT(refreshToken.UserID, c.JwtSecret, accessTokenExpires)
	if err != nil {
		respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error: failed to create new access token"})
		return
	}

	// If all good, return the new access token
	respondJson(w, http.StatusOK, validResponse{Token: newToken})
}

// handlerRevoke is an HTTP handler function to revoke the refresh token
// The endpoint does not accept a request body, but expects the refresh token to be sent in the Authorization header
func (c *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Get Bearer bearer from the request headers
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJson(w, http.StatusUnauthorized, errorResponse{Error: "No authorization token provided"})
		return
	}

	// Check if refresh token is valid
	refreshToken, err := c.db.GetRefreshTokenByToken(r.Context(), bearer)
	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		respondJson(w, http.StatusUnauthorized, errorResponse{Error: "Invalid or expired refresh token"})
		return
	}

	// Revoke the refresh token by setting revoked_at and updated_at to now (native SQL)
	err = c.db.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondJson(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error: failed to revoke refresh token"})
		return
	}

	// If all good, just return 204 with no body
	w.WriteHeader(http.StatusNoContent)
}
