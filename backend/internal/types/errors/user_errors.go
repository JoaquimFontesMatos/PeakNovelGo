package errors

import (
	"backend/internal/types"
	"net/http"
)

const (
	// Validation errors
	INVALID_PASSWORD            = "INVALID_PASSWORD"
	INVALID_EMAIL               = "INVALID_EMAIL"
	INVALID_USERNAME            = "INVALID_USERNAME"
	INVALID_BIO                 = "INVALID_BIO"
	INVALID_PROFILE_PICTURE     = "INVALID_PROFILE_PICTURE"
	INVALID_PREFERRED_LANGUAGE  = "INVALID_PREFERRED_LANGUAGE"
	INVALID_READING_PREFERENCES = "INVALID_READING_PREFERENCES"
	INVALID_ROLES               = "INVALID_ROLES"
	INVALID_DATE_OF_BIRTH       = "INVALID_DATE_OF_BIRTH"

	// User errors
	CREATING_USER_ERROR    = "CREATING_USER_ERROR"
	UPDATING_USER_ERROR    = "UPDATING_USER_ERROR"
	UPDATING_EMAIL_ERROR   = "UPDATING_EMAIL_ERROR"
	UPDATING_FIELDS_ERROR  = "UPDATING_FIELDS_ERROR"
	USER_ALREADY_EXISTS    = "USER_ALREADY_EXISTS"
	USER_NOT_FOUND         = "USER_NOT_FOUND"
	USER_DEACTIVATED_ERROR = "USER_DEACTIVATED_ERROR"
)

var (
	// Password errors
	ErrRequiredPassword = &types.MyCustomError{
		Message:    "Password is required",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_PASSWORD,
	}
	ErrShortPassword = &types.MyCustomError{
		Message:    "Password must be at least 8 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_PASSWORD,
	}
	ErrLongPassword = &types.MyCustomError{
		Message:    "Password cannot be longer than 72 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_PASSWORD,
	}

	// Email errors
	ErrEmailRequired = &types.MyCustomError{
		Message:    "Email is required",
		StatusCode: http.StatusBadRequest,
		Code:       "INVALID_EMAIL",
	}
	ErrEmailTooLong = &types.MyCustomError{
		Message:    "Email cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_EMAIL,
	}
	ErrEmailTooShort = &types.MyCustomError{
		Message:    "Email must be at least 1 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_EMAIL,
	}
	ErrInvalidEmailFormat = &types.MyCustomError{
		Message:    "Invalid email format",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_EMAIL,
	}

	// Username errors
	ErrUsernameRequired = &types.MyCustomError{
		Message:    "Username is required",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_USERNAME,
	}
	ErrUsernameTooLong = &types.MyCustomError{
		Message:    "Username cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_USERNAME,
	}
	ErrUsernameTooShort = &types.MyCustomError{
		Message:    "Username must be at least 1 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_USERNAME,
	}

	// Bio errors
	ErrBioTooLong = &types.MyCustomError{
		Message:    "Bio cannot be longer than 500 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_BIO,
	}

	// Profile picture errors
	ErrProfilePictureTooLong = &types.MyCustomError{
		Message:    "Profile picture URL cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_PROFILE_PICTURE,
	}

	// Preferred language errors
	ErrPreferredLanguageTooLong = &types.MyCustomError{
		Message:    "Preferred language cannot be longer than 100 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_PREFERRED_LANGUAGE,
	}

	// Reading preferences errors
	ErrReadingPreferencesTooLong = &types.MyCustomError{
		Message:    "Reading preferences cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_READING_PREFERENCES,
	}

	// Roles errors
	ErrRolesTooLong = &types.MyCustomError{
		Message:    "Roles cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_ROLES,
	}

	// Date of birth errors
	ErrBirthDateTooYoung = &types.MyCustomError{
		Message:    "You must be at least 18 years old",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_DATE_OF_BIRTH,
	}

	// User errors
	ErrUserDeactivatedError = &types.MyCustomError{
		Message:    "User account is deactivated",
		StatusCode: http.StatusForbidden,
		Code:       USER_DEACTIVATED_ERROR,
	}
	ErrUserNotFound = &types.MyCustomError{
		Message:    "User not found",
		StatusCode: http.StatusNotFound,
		Code:       USER_NOT_FOUND,
	}
	ErrUserAlreadyExists = &types.MyCustomError{
		Message:    "User already exists",
		StatusCode: http.StatusConflict,
		Code:       USER_ALREADY_EXISTS,
	}
)
