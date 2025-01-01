package validators

type UserError struct {
    Code    string
    Message string
}

func (e *UserError) Error() string {
    return e.Message
}

var (
    ErrUserNotFound    = &UserError{Code: "USER_NOT_FOUND", Message: "User not found"}
    ErrUserDeleted     = &UserError{Code: "USER_DELETED", Message: "User account is deactivated"}
    ErrInvalidPassword = &UserError{Code: "INVALID_PASSWORD", Message: "Invalid password"}
    ErrPasswordDiff    = &UserError{Code: "PASSWORD_DIFF", Message: "New password cannot be the same as the current password"}
    ErrTokenExpired    = &UserError{Code: "INVALID_TOKEN", Message: "Invalid token or token expired"}
)
