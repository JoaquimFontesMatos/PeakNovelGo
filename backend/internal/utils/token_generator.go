package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
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

func ExtractToken(headerValue string) (string, error) {
    if !strings.HasPrefix(headerValue, "Bearer ") {
        return "", fmt.Errorf("invalid token format")
    }
    return strings.TrimPrefix(headerValue, "Bearer "), nil
}
