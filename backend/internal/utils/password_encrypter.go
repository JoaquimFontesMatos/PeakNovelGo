package utils

import (
	"backend/internal/types"
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
//   - INTERNAL_ERROR if there is an error hashing the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to hash password", err)
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
