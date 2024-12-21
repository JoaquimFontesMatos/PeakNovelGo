package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// GenerateVerificationToken generates a secure random verification token.
func GenerateVerificationToken() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("Failed to generate random bytes for token: %v", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
