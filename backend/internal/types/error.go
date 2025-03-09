package types

// HTTPError is an interface that defines the common methods for an HTTP error.
type HTTPError interface {
	error
	HTTPStatus() int
	ErrorCode() string
}

// MyCustomError MyError is a custom error type that implements the HTTPError interface.
type MyCustomError struct {
	Message    string
	StatusCode int
	Code       string
	Wrapped    error
}

// Error returns the error message.
//
// Returns:
//   - string (The error message.)
func (e *MyCustomError) Error() string { return e.Message }

// HTTPStatus returns the HTTP status code.
//
// Returns:
//   - int (The HTTP status code.)
func (e *MyCustomError) HTTPStatus() int { return e.StatusCode }

// ErrorCode returns the error code.
//
// Returns:
//   - string (The error code.)
func (e *MyCustomError) ErrorCode() string { return e.Code }

// WrappedError returns the wrapped error.
//
// Returns:
//   - error (The wrapped error.)
func (e *MyCustomError) WrappedError() error { return e.Wrapped }

// WrapError wraps an error with a custom error type.
//
// Parameters:
//   - code (string): Error code.
//   - message (string): Error message.
//   - statusCode (int): HTTP status code.
//   - err (error): The original error being wrapped.
//
// Returns:
//   - *MyCustomError: A pointer to the wrapped error.
func WrapError(code string, message string, statusCode int, err error) *MyCustomError {
	return &MyCustomError{
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
		Wrapped:    err,
	}
}
