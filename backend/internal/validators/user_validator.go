package validators

import (
	"backend/internal/dtos"
	"backend/internal/models"
	"backend/internal/types"
	"backend/internal/types/errors"
	"net/http"
	"regexp"
	"time"
)

// ValidateUser validates a user model.
// It checks the email, password, and other user fields for validity.
//
// Parameters:
//   - user (*models.User): The user model to validate.
//
// Returns:
//   - error: An error if any validation fails, nil otherwise.
//
// Error types:
//   - wrapped http.StatusBadRequest errors
func ValidateUser(user *models.User) error {
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}

	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	fields := dtos.UpdateRequest{
		Bio:                user.Bio,
		PreferredLanguage:  user.PreferredLanguage,
		Roles:              user.Roles,
		Username:           user.Username,
		ProfilePicture:     user.ProfilePicture,
		ReadingPreferences: user.ReadingPreferences,
		DateOfBirth:        user.DateOfBirth.Format("2006-01-02"),
	}

	if err := ValidateUserFields(fields); err != nil {
		return err
	}

	return nil
}

// ValidatePassword checks if the provided password meets certain criteria.
//
// Parameters:
//   - password (string): The password to validate.
//
// Returns:
//   - error: nil if the password is valid, otherwise an error indicating the reason for invalidity.
//
// Error types:
//   - errors.ErrRequiredPassword: Returned if the password is empty.
//   - errors.ErrShortPassword: Returned if the password is less than 8 characters long.
//   - errors.ErrLongPassword: Returned if the password is longer than 72 characters.
func ValidatePassword(password string) error {
	if password == "" {
		return errors.ErrRequiredPassword
	}
	if len(password) < 8 {
		return errors.ErrShortPassword
	}
	if len(password) > 72 {
		return errors.ErrLongPassword
	}
	return nil
}

// ValidateEmail checks if the provided email is valid.
//
// Parameters:
//   - email (string): The email address to validate.
//
// Returns:
//   - error: An error if the email is invalid, nil otherwise.
//
// Error types:
//   - errors.ErrEmailRequired: If the email is empty.
//   - errors.ErrEmailTooLong: If the email is longer than 255 characters.
//   - errors.ErrInvalidEmailFormat: If the email format is invalid.
//   - types.WrapError(errors.INVALID_EMAIL, ..., http.StatusBadRequest, ...): If an error occurred during regex matching.
func ValidateEmail(email string) error {
	if email == "" {
		return errors.ErrEmailRequired
	}

	if len(email) > 255 {
		return errors.ErrEmailTooLong
	}

	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return types.WrapError(errors.INVALID_EMAIL, "Failed to validate email format", http.StatusBadRequest, err)
	}

	if !matched {
		return errors.ErrInvalidEmailFormat
	}
	return nil
}

// ValidateUsername checks if the provided username is valid.
//
// Parameters:
//   - username (string): The username to validate.
//
// Returns:
//   - error: An error if the username is invalid, nil otherwise.
//
// Error types:
//   - errors.ErrUsernameRequired: Returned if the username is empty.
//   - errors.ErrUsernameTooLong: Returned if the username is longer than 255 characters.
func ValidateUsername(username string) error {
	if username == "" {
		return errors.ErrUsernameRequired
	}
	if len(username) > 255 {
		return errors.ErrUsernameTooLong
	}
	return nil
}

// ValidateUserFields validates the fields of a user update request.
// It checks the length of various string fields and validates the date of birth.
//
// Parameters:
//   - fields (dtos.UpdateRequest): The user update request fields.
//
// Returns:
//   - error: An error if any of the validations fail, nil otherwise.
//
// Error types:
//   - errors.ErrUsernameTooLong: If the username is too long.
//   - errors.ErrBioTooLong: If the bio is too long.
//   - errors.ErrProfilePictureTooLong: If the profile picture URL is too long.
//   - errors.ErrPreferredLanguageTooLong: If the preferred language is too long.
//   - errors.ErrReadingPreferencesTooLong: If the reading preferences are too long.
//   - errors.ErrRolesTooLong: If the roles string is too long.
//   - errors.INVALID_DATE_OF_BIRTH: If the date of birth is not a valid date or is in the wrong format.
//   - errors.ErrBirthDateTooYoung: If the user is younger than 18 based on provided date of birth.
func ValidateUserFields(fields dtos.UpdateRequest) error {
	// Check the username length
	if len(fields.Username) > 255 {
		return errors.ErrUsernameTooLong
	}

	// Check the bio length
	if len(fields.Bio) > 500 {
		return errors.ErrBioTooLong
	}

	// Check the profile picture URL length
	if len(fields.ProfilePicture) > 255 {
		return errors.ErrProfilePictureTooLong
	}

	// Check the preferred language length
	if len(fields.PreferredLanguage) > 100 {
		return errors.ErrPreferredLanguageTooLong
	}

	// Check the reading preferences length
	if len(fields.ReadingPreferences) > 255 {
		return errors.ErrReadingPreferencesTooLong
	}

	// Check the roles length
	if len(fields.Roles) > 255 {
		return errors.ErrRolesTooLong
	}

	// Validate the date of birth
	if fields.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", fields.DateOfBirth)
		if err != nil {
			return types.WrapError(errors.INVALID_DATE_OF_BIRTH, "Date of birth must be a valid date in YYYY-MM-DD format", http.StatusBadRequest, err)
		}
		if dob.After(time.Now().AddDate(-18, 0, 0)) {
			return errors.ErrBirthDateTooYoung
		}
	}

	return nil
}

// ValidateIsDeleted checks if the user is deleted.
//
// Parameters:
//   - user (models.User): The user to check.
//
// Returns:
//   - error: errors.ErrUserDeactivatedError if the user is deleted, nil otherwise.
func ValidateIsDeleted(user models.User) error {
	if user.IsDeleted {
		return errors.ErrUserDeactivatedError
	}

	return nil
}

// ValidateIsNewPasswordTheSame checks if the new password is different from the current password.
//
// Parameters:
//   - currentPassword (string): The current user's password.
//   - newPassword (string): The new password entered by the user.
//
// Returns:
//   - error: errors.ErrPasswordDiff if the new password is the same as the current password, nil otherwise.
func ValidateIsNewPasswordTheSame(currentPassword string, newPassword string) error {
	if currentPassword == newPassword {
		return errors.ErrPasswordDiff
	}
	return nil
}
