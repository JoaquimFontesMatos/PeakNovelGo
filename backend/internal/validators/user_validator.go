package validators

import (
	"backend/internal/dtos"
	"backend/internal/models"
	"backend/internal/types"
	"regexp"
	"time"
)

// ValidateUser validates the user input and returns an error if the user is invalid.
//
// Parameters:
//   - user (*models.User): A pointer to the User struct containing user information.
//
// Returns:
//   - VALIDATION_ERROR if the user input is invalid.
//   - USER_DEACTIVATED if the user's account is deactivated.
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

// ValidatePassword validates if the password is valid
//
// Parameters:
//   - password string (password)
//
// Returns:
//   - VALIDATION_ERROR if the password is invalid
func ValidatePassword(password string) error {
	if password == "" {
		return types.WrapError("VALIDATION_ERROR", "Password is required", nil)
	}
	if len(password) < 8 {
		return types.WrapError("VALIDATION_ERROR", "Password must be at least 8 characters long", nil)
	}
	if len(password) > 72 {
		return types.WrapError("VALIDATION_ERROR", "Password cannot be longer than 72 characters", nil)
	}
	return nil
}

// ValidateEmail validates if the email is valid
//
// Parameters:
//   - email string (email)
//
// Returns:
//   - VALIDATION_ERROR if the email is invalid
func ValidateEmail(email string) error {
	if email == "" {
		return types.WrapError("VALIDATION_ERROR", "Email is required", nil)
	}

	if len(email) > 255 {
		return types.WrapError("VALIDATION_ERROR", "Email cannot be longer than 255 characters", nil)
	}

	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return types.WrapError("VALIDATION_ERROR", "Failed to validate email format", err)
	}
	if !matched {
		return types.WrapError("VALIDATION_ERROR", "Invalid email format", nil)
	}
	return nil
}

// ValidateUsername validates if the username is valid
//
// Parameters:
//   - username string (username)
//
// Returns:
//   - VALIDATION_ERROR if the username is invalid
func ValidateUsername(username string) error {
	if username == "" {
		return types.WrapError("VALIDATION_ERROR", "Username is required", nil)
	}
	if len(username) > 255 {
		return types.WrapError("VALIDATION_ERROR", "Username cannot be longer than 255 characters", nil)
	}
	return nil
}

// ValidateUserFields validates the fields being updated and ensures they meet validation criteria.
// ValidateUserFields validates an UpdateFields struct.
//
// Parameters:
//   - fields dtos.UpdateRequest (UpdateRequest struct)
//
// Returns:
//   - VALIDATION_ERROR if the fields are invalid
func ValidateUserFields(fields dtos.UpdateRequest) error {
	// Check the username length
	if len(fields.Username) > 255 {
		return types.WrapError("VALIDATION_ERROR", "Username cannot be longer than 255 characters", nil)
	}

	// Check the bio length
	if len(fields.Bio) > 500 {
		return types.WrapError("VALIDATION_ERROR", "Bio cannot be longer than 500 characters", nil)
	}

	// Check the profile picture URL length
	if len(fields.ProfilePicture) > 255 {
		return types.WrapError("VALIDATION_ERROR", "Profile picture URL cannot be longer than 255 characters", nil)
	}

	// Check the preferred language length
	if len(fields.PreferredLanguage) > 100 {
		return types.WrapError("VALIDATION_ERROR", "Preferred language cannot be longer than 100 characters", nil)
	}

	// Check the reading preferences length
	if len(fields.ReadingPreferences) > 255 {
		return types.WrapError("VALIDATION_ERROR", "Reading preferences cannot be longer than 255 characters", nil)
	}

	// Check the roles length
	if len(fields.Roles) > 255 {
		return types.WrapError("VALIDATION_ERROR", "Roles cannot be longer than 255 characters", nil)
	}

	// Validate the date of birth
	if fields.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", fields.DateOfBirth)
		if err != nil {
			return types.WrapError("VALIDATION_ERROR", "Date of birth must be a valid date in YYYY-MM-DD format", err)
		}
		if dob.After(time.Now().AddDate(-18, 0, 0)) {
			return types.WrapError("VALIDATION_ERROR", "You must be at least 18 years old", nil)
		}
	}

	return nil
}

// ValidateIsDeleted validates if the user is deactivated
//
// Parameters:
//   - user models.User (User struct)
//
// Returns:
//   - ErrUserDeactivated if the user is deactivated
func ValidateIsDeleted(user models.User) error {
	if user.IsDeleted {
		return types.ErrUserDeactivated
	}

	return nil
}

// ValidateIsNewPasswordTheSame validates if the new password is the same as the current password
//
// Parameters:
//   - currentPassword string (current password)
//   - newPassword string (new password)
//
// Returns:
//   - ErrPasswordDiff if the new password is the same as the current password
func ValidateIsNewPasswordTheSame(currentPassword string, newPassword string) error {
	if currentPassword == newPassword {
		return types.ErrPasswordDiff
	}
	return nil
}
