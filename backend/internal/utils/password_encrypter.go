package utils

import (
	"backend/internal/types"
	"backend/internal/types/errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const MaxPasswordLength = 72

// HashPassword generates a bcrypt hash of the password.
//
// Parameters:
//   - password string (password)
//
// Returns:
//   - string (hashed password)
//   - error (error with status code: http.StatusInternalServerError if the password could not be encrypted)
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", types.WrapError(errors.HASH_PASSWORD, "Failed to hash password", http.StatusInternalServerError, err)
	}
	return string(hash), nil
}

// ComparePassword compares a hashed password with its plaintext version.
//
// Parameters:
//   - hashedPassword string (hashed password)
//   - plainPassword string (plaintext password)
//
// Returns:
//   - bool (true if the passwords match, false otherwise)
func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
