package errors

import (
	"backend/internal/types"
	"net/http"
)

const (
	DB_UNAVAILABLE = "DB_UNAVAILABLE"
	DB_TIMEOUT     = "DB_TIMEOUT"
	// Add more database error codes here
)

var ErrDatabaseOffline = &types.MyCustomError{
	Message:    "Database unavailable",
	StatusCode: http.StatusServiceUnavailable,
	Code:       DB_UNAVAILABLE,
}
