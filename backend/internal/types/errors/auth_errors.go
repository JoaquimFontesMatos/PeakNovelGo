package errors

import (
	"backend/internal/types"
	"net/http"
)

const (
	GENERATING_TOKEN      = "GENERATING_TOKEN"
	INVALID_TOKEN         = "INVALID_TOKEN"
	REVOKING_TOKEN        = "REVOKING_TOKEN"
	REFRESH_TOKEN_REVOKED = "REFRESH_TOKEN_REVOKED"
	PASSWORD_DIFF         = "PASSWORD_DIFF"
	HASH_PASSWORD         = "HASH_PASSWORD"
	INVALID_CREDENTIALS   = "INVALID_CREDENTIALS"
)

var (
	// Token errors
	ErrInvalidToken = &types.MyCustomError{
		Message:    "Invalid token (cannot be longer than 255 characters)",
		StatusCode: http.StatusUnauthorized,
		Code:       INVALID_TOKEN,
	}
	ErrRefreshTokenRevoked = &types.MyCustomError{
		Message:    "Refresh token has been revoked",
		StatusCode: http.StatusUnauthorized,
		Code:       REFRESH_TOKEN_REVOKED,
	}

	// Password errors
	ErrPasswordDiff = &types.MyCustomError{
		Message:    "New password must be different from the current password",
		StatusCode: http.StatusUnauthorized,
		Code:       PASSWORD_DIFF,
	}
	ErrPasswordDoesNotMatch = &types.MyCustomError{
		Message:    "New password does not match the current password",
		StatusCode: http.StatusUnauthorized,
		Code:       PASSWORD_DIFF,
	}
	ErrInvalidCredentials = &types.MyCustomError{
		Message:    "Invalid credentials",
		StatusCode: http.StatusUnauthorized,
		Code:       INVALID_CREDENTIALS,
	}
)
