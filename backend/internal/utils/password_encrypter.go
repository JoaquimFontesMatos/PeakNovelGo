package utils

import (
	"backend/internal/types"
	"golang.org/x/crypto/bcrypt"
)

const MaxPasswordLength = 72

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", types.WrapError("INTERNAL_SERVER_ERROR", "Failed to hash password", err)
	}
	return string(hash), nil
}

// ComparePassword compares a hashed password with its plaintext version.
func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
