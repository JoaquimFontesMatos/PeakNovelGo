package validators

import (
	"backend/internal/types/errors"

	"time"
)

// IsVerificationTokenExpired checks if the verification token is expired.
//
// Parameters:
//   - createdAt (time.Time): time when the token was created
//   - emailVerified (bool):whether the email has been verified
//
// Returns:
//   - bool: true if the token is expired, false otherwise
func IsVerificationTokenExpired(createdAt time.Time, emailVerified bool) bool {
	timeDifference := time.Since(createdAt)

	// Check if the token is expired
	if timeDifference >= 60*time.Minute || emailVerified {
		return true
	}
	return false
}

// ValidateToken checks if a token is valid.
//
// Parameters:
//   - token (string): The token to validate.
//
// Returns:
//   - error: An error (errors.ErrInvalidToken) if the token is invalid, nil otherwise.
func ValidateToken(token string) error {
	if token == "" {
		return errors.ErrInvalidToken
	}

	if len(token) > 255 {
		return errors.ErrInvalidToken
	}

	return nil
}
