package types

import "net/http"

type HTTPError interface {
	error
	HTTPStatus() int
	ErrorCode() string
}

// Example implementation
type MyCustomError struct {
	Message    string
	StatusCode int
	Code       string
	Wrapped    error
}

func (e *MyCustomError) Error() string     { return e.Message }
func (e *MyCustomError) HTTPStatus() int   { return e.StatusCode }
func (e *MyCustomError) ErrorCode() string { return e.Code }

var ErrDatabaseOffline = &MyCustomError{
    Message:    "Database unavailable",
    StatusCode: http.StatusServiceUnavailable,
    Code:       "DB_UNAVAILABLE",
}