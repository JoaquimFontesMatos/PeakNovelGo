package repositories

import (
	"errors"
	"log"

	"backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) CreateUser(user *models.User) error {
	// Create a new user in the database
	if err := r.db.Create(user).Error; err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	user := &models.User{}
	// Fetch user by ID
	if err := r.db.First(user, id).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByVerificationToken(token string) (*models.User, error) {
	user := &models.User{}
	// Fetch user by verification token
	if err := r.db.Where("verification_token = ?", token).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	// Update user in the database
	if err := r.db.Save(user).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(user *models.User) error {
	// Soft delete the user in the database
	user.IsDeleted = true

	return r.UpdateUser(user)
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	// Fetch user by email
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	// Fetch user by username
	if err := r.db.Where("username = ?", username).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUserEmail(userID uint, newEmail string, token string) error {
	// Update user email and verification token
	result := r.db.Model(&models.User{}).Where("id = ? AND is_deleted = 0", userID).
		Updates(map[string]interface{}{
			"email":              newEmail,
			"verification_token": token,
			"email_verified":     false,
		})

	if result.RowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}

	if result.Error != nil {
		log.Printf("Failed to update user email: %v", result.Error)
		return result.Error
	}

	return nil
}


func (r *UserRepository) UpdateUserFields(userID uint, fields interface{}) error {
	// Update user fields
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Updates(fields).Error; err != nil {
		log.Printf("Failed to update user fields: %v", err)
		return err
	}
	return nil
}
