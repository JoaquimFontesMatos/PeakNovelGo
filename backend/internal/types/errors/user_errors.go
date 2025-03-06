package errors

import (
	"backend/internal/types"
	"net/http"
)

var (
	ErrInvalidRole = &types.MyCustomError{
		Message:    "Invalid user role",
		StatusCode: http.StatusForbidden,
		Code:       "INVALID_ROLE",
	}
)