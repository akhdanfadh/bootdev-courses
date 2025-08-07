package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// MakeJWT creates a new JWT token for the given user ID, signing it with the provided secret and setting an expiration time
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Claims are statements about the entity to be encoded in the JWT
	// registered claims are standard claims defined by the JWT specification
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	// Create a new JWT token with the claims and HMAC signing method
	// https://golang-jwt.github.io/jwt/usage/signing_methods/#signing-methods-and-key-types
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // not yet signed
	ss, err := token.SignedString([]byte(tokenSecret))         // final signed token
	if err != nil {
		return "", err
	}
	return ss, nil
}

// ValidateJWT validates the given JWT token string using the provided secret and returns the user ID if valid
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Necessary arguments to parse the token
	claims := &jwt.RegisteredClaims{} // to hold the parsed claims
	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (any, error) {
		// to tell the parser what key to use for validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil // return the appropriate key
	}

	// Parse the token
	_, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		// this will catch expiration, invalid signature, malformed token, etc.
		return uuid.Nil, fmt.Errorf("token validation failed: %w", err)
	}
	if claims.Issuer != "chirpy" { // but we still need to check issuer (us)
		return uuid.Nil, fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}

	// If all good, return the user ID from the claims
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %w", err)
	}
	return userID, nil
}
