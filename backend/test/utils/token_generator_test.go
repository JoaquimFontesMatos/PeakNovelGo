package utils_test

import (
	"backend/internal/utils"
	"testing"
)

func TestGenerateVerificationToken(t *testing.T) {
	token := utils.GenerateVerificationToken()
	if len(token) == 0 {
		t.Errorf("Expected a non-empty token, but got an empty string")
	}
}

func TestGenerateVerificationToken_Length(t *testing.T) {
	token := utils.GenerateVerificationToken()

	// Length will be 44 for 32 bytes encoded in Base64
	expectedLength := 44
	if len(token) != expectedLength {
		t.Errorf("Expected token length to be %d, but got %d", expectedLength, len(token))
	}
}

func TestGenerateVerificationToken_Uniqueness(t *testing.T) {
	token1 := utils.GenerateVerificationToken()
	token2 := utils.GenerateVerificationToken()

	if token1 == token2 {
		t.Errorf("Expected tokens to be unique, but got identical values: %s", token1)
	}
}
