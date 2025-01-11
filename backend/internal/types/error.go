package types

// MyError struct represents a custom error with a code and message.
type MyError struct {
	Code    string
	Message string
	Err     error
}

// Error returns the error message.
//
// Returns:
//   - string (error message)
func (e *MyError) Error() string {
	return e.Message
}

// WrapError wraps an error with a code and message.
//
// Parameters:
//   - code string (error code)
//   - message string (error message)
//   - err error (error to wrap)
//
// Returns:
//   - *MyError (MyError struct)
func WrapError(code string, message string, err error) *MyError {
	return &MyError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// static error codes
var (
	ErrUserNotFound    = &MyError{Code: "USER_NOT_FOUND_ERROR", Message: "User not found"}
	ErrInvalidPassword = &MyError{Code: "INVALID_PASSWORD_ERROR", Message: "Invalid password"}
	ErrPasswordDiff    = &MyError{Code: "PASSWORD_DIFF_ERROR", Message: "New password cannot be the same as the current password"}
	ErrTokenExpired    = &MyError{Code: "INVALID_TOKEN_ERROR", Message: "Invalid token or token expired"}
	ErrUserDeactivated = &MyError{Code: "USER_DEACTIVATED_ERROR", Message: "User account is deactivated"}
)
