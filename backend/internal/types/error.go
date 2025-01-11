package types

// MyError struct represents a custom error with a code and message.
type MyError struct {
	Code    ErrorType
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
func WrapError(code ErrorType, message string, err error) *MyError {
	return &MyError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

type ErrorType string

const (
	VALIDATION_ERROR            ErrorType = "VALIDATION_ERROR"
	REPOSITORY_ERROR            ErrorType = "REPOSITORY_ERROR"
	SERVICE_ERROR               ErrorType = "SERVICE_ERROR"
	AUTHENTICATION_ERROR        ErrorType = "AUTHENTICATION_ERROR"
	AUTHORIZATION_ERROR         ErrorType = "AUTHORIZATION_ERROR"
	INTERNAL_SERVER_ERROR       ErrorType = "INTERNAL_SERVER_ERROR"
	USER_DEACTIVATED_ERROR      ErrorType = "USER_DEACTIVATED_ERROR"
	INVALID_TOKEN_ERROR         ErrorType = "INVALID_TOKEN_ERROR"
	INVALID_PASSWORD_ERROR      ErrorType = "INVALID_PASSWORD_ERROR"
	PASSWORD_DIFF_ERROR         ErrorType = "PASSWORD_DIFF_ERROR"
	USER_NOT_FOUND_ERROR        ErrorType = "USER_NOT_FOUND_ERROR"
	EMAIL_SEND_ERROR            ErrorType = "EMAIL_SEND_ERROR"
	INVALID_BODY_ERROR          ErrorType = "INVALID_BODY_ERROR"
	INVALID_ID_ERROR            ErrorType = "INVALID_ID_ERROR"
	CONFLICT_ERROR              ErrorType = "CONFLICT_ERROR"
	REFRESH_TOKEN_REVOKED_ERROR ErrorType = "REFRESH_TOKEN_REVOKED_ERROR"
	INVALID_CREDENTIALS_ERROR   ErrorType = "INVALID_CREDENTIALS_ERROR"
	// Add more types as needed
)
