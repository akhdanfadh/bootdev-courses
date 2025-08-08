package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// GetBearerToken extracts the Bearer token from the Authorization header
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:], nil
	}
	return "", fmt.Errorf("missing bearer token in Authorization header")
}

// GetApiKey extracts the API key from the Authorization header
func GetApiKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if strings.HasPrefix(apiKey, "ApiKey ") {
		return apiKey[7:], nil
	}
	return "", fmt.Errorf("missing ApiKey in Authorization header")
}
