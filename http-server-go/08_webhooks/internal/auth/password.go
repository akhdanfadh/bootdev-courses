package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the given password using bcrypt and returns the hashed password as a string
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPasswordHash compares a password with a hashed password and returns an error if they do not match
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
