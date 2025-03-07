package types

// HTTPError is an interface that defines the common methods for an HTTP error.
type HTTPError interface {
	error
	HTTPStatus() int
	ErrorCode() string
}

// MyError is a custom error type that implements the HTTPError interface.
type MyCustomError struct {
	Message    string
	StatusCode int
	Code       string
	Wrapped    error
}

// Error returns the error message.
func (e *MyCustomError) Error() string     { return e.Message }

// HTTPStatus returns the HTTP status code.
func (e *MyCustomError) HTTPStatus() int   { return e.StatusCode }

// ErrorCode returns the error code.
func (e *MyCustomError) ErrorCode() string { return e.Code }

// WrapError wraps an error with a code and message.
//
// Parameters:
//   - code string (error code)
//   - message string (error message)
//   - statusCode int (HTTP status code)
//   - err error (error to wrap)
//
// Returns:
//   - *MyError (MyError struct)
func WrapError(code string, message string, statusCode int, err error) *MyCustomError {
	return &MyCustomError{
		Code:    code,
		StatusCode: statusCode,
		Message: message,
		Wrapped:     err,
	}
}
