package utils

import (
	"backend/internal/types/errors"
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"
)

// GenerateVerificationToken generates a secure random verification token.
//
// Returns:
//   - string (verification token)
func GenerateVerificationToken() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("Failed to generate random bytes for token: %v", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// ExtractToken extracts the token from the Authorization header.
//
// Parameters:
//   - headerValue string (Authorization header value)
//
// Returns:
//   - string (token)
//   - INVALID_TOKEN_ERROR if the token is invalid
func ExtractToken(headerValue string) (string, error) {
	if !strings.HasPrefix(headerValue, "Bearer ") {
		return "", errors.ErrInvalidToken
	}
	return strings.TrimPrefix(headerValue, "Bearer "), nil
}
