package repositories

import (
	"log"

	"backend/internal/models"
	"backend/internal/types"
	"backend/internal/validators"

	"gorm.io/gorm"
)

// UserRepository struct represents a repository for user management.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
//
// Parameters:
//   - db *gorm.DB (Gorm database connection)
//
// Returns:
//   - *UserRepository (pointer to UserRepository struct)
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database.
//
// Parameters:
//   - user *models.User (pointer to User struct)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if the user could not be inserted
func (r *UserRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		log.Printf("Failed to insert user: %v", err)
		return types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to insert user", err)
	}
	return nil
}

// GetUserByID gets a user by ID.
//
// Parameters:
//   - id uint (ID of the user)
//
// Returns:
//   - *models.User (pointer to User struct)
//   - INTERNAL_SERVER_ERROR if the user could not be fetched
//   - USER_DEACTIVATED_ERROR if the user is deactivated
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	user := &models.User{}

	if err := r.db.First(user, id).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, types.WrapError(types.USER_NOT_FOUND_ERROR, "User not found", err)
	}

	if err := validators.ValidateIsDeleted(*user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByVerificationToken gets a user by verification token.
//
// Parameters:
//   - token string (verification token)
//
// Returns:
//   - *models.User (pointer to User struct)
//   - INTERNAL_SERVER_ERROR if the user could not be fetched
//   - USER_DEACTIVATED_ERROR if the user is deactivated
func (r *UserRepository) GetUserByVerificationToken(token string) (*models.User, error) {
	user := &models.User{}

	if err := r.db.Where("verification_token = ?", token).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, types.WrapError(types.USER_NOT_FOUND_ERROR, "User not found", err)
	}

	if err := validators.ValidateIsDeleted(*user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user in the database.
//
// Parameters:
//   - user *models.User (pointer to User struct)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if the user could not be updated
func (r *UserRepository) UpdateUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		return types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to update user", err)
	}

	return nil
}

// DeleteUser deactivates a user in the database.
//
// Parameters:
//   - user *models.User (pointer to User struct)
//
// Returns:
//   - error (INTERNAL_SERVER_ERROR if the user could not be deactivated)
func (r *UserRepository) DeleteUser(user *models.User) error {
	user.IsDeleted = true

	return r.UpdateUser(user)
}

// GetUserByEmail gets a user by email.
//
// Parameters:
//   - email string (email of the user)
//
// Returns:
//   - *models.User (pointer to User struct)
//   - INTERNAL_SERVER_ERROR if the user could not be fetched
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, types.WrapError(types.USER_NOT_FOUND_ERROR, "User not found", err)
	}

	if err := validators.ValidateIsDeleted(*user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername gets a user by username.
//
// Parameters:
//   - username string (username of the user)
//
// Returns:
//   - *models.User (pointer to User struct)	º
//   - INTERNAL_SERVER_ERROR if the user could not be fetched
//   - USER_DEACTIVATED_ERROR if the user is deactivated
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}

	if err := r.db.Where("username = ?", username).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)

		return nil, types.WrapError(types.USER_NOT_FOUND_ERROR, "User not found", err)
	}

	if err := validators.ValidateIsDeleted(*user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserEmail updates the email and verification token of a user.
//
// Parameters:
//   - userID uint (ID of the user)
//   - newEmail string (new email of the user)
//   - token string (new verification token)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if the user was not found, deactivated or the update failed
func (r *UserRepository) UpdateUserEmail(userID uint, newEmail string, token string) error {
	result := r.db.Model(&models.User{}).Where("id = ? AND is_deleted = 0", userID).
		Updates(map[string]interface{}{
			"email":              newEmail,
			"verification_token": token,
			"email_verified":     false,
		})

	if result.RowsAffected == 0 {
		return types.WrapError(types.INTERNAL_SERVER_ERROR, "User not found or already deleted", nil)
	}

	if result.Error != nil {
		log.Printf("Failed to update user email: %v", result.Error)
		return types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to update user email", result.Error)
	}

	return nil
}

// UpdateUserFields updates the fields of a user.
//
// Parameters:
//   - userID uint (ID of the user)
//   - fields interface{} (fields to update)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if the update failed
func (r *UserRepository) UpdateUserFields(userID uint, fields interface{}) error {
	// Update user fields
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Updates(fields).Error; err != nil {
		log.Printf("Failed to update user fields: %v", err)
		return types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to update user fields", err)
	}
	return nil
}
