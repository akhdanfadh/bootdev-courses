package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	// Create some JWT tokens for testing
	userID := uuid.New()
	secret := "test-secret"
	validToken, _ := MakeJWT(userID, secret, time.Hour)
	expiredToken, _ := MakeJWT(userID, secret, -time.Hour)

	// Create token with wrong issuer
	wrongIssuerClaims := &jwt.RegisteredClaims{
		Issuer:    "wrong-issuer",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   userID.String(),
	}
	wrongIssuerTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, wrongIssuerClaims)
	wrongIssuerToken, _ := wrongIssuerTokenObj.SignedString([]byte(secret))

	// Create token with invalid subject
	invalidSubjectClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   "not-a-valid-uuid",
	}
	invalidSubjectTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidSubjectClaims)
	invalidSubjectToken, _ := invalidSubjectTokenObj.SignedString([]byte(secret))

	tests := map[string]struct {
		tokenString string
		tokenSecret string
		expectedID  uuid.UUID
		wantErr     bool
	}{
		"Valid token": {
			tokenString: validToken,
			tokenSecret: secret,
			expectedID:  userID,
			wantErr:     false,
		},
		"Expired token": {
			tokenString: expiredToken,
			tokenSecret: secret,
			expectedID:  uuid.Nil,
			wantErr:     true,
		},
		"Wrong secret": {
			tokenString: validToken,
			tokenSecret: "wrong-secret",
			expectedID:  uuid.Nil,
			wantErr:     true,
		},
		"Malformed token": {
			tokenString: "not.a.valid.token",
			tokenSecret: secret,
			expectedID:  uuid.Nil,
			wantErr:     true,
		},
		"Empty token": {
			tokenString: "",
			tokenSecret: secret,
			expectedID:  uuid.Nil,
			wantErr:     true,
		},
		"Wrong issuer": {
			tokenString: wrongIssuerToken,
			tokenSecret: secret,
			expectedID:  uuid.Nil,
			wantErr:     true,
		},
		"Invalid subject UUID": {
			tokenString: invalidSubjectToken,
			tokenSecret: secret,
			expectedID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.expectedID {
				t.Errorf("ValidateJWT() gotID = %v, expectedID %v", gotID, tt.expectedID)
			}
		})
	}
}
