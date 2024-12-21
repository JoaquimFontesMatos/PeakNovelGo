package repositories

import (
	"log"

	"backend/internal/models"
	"backend/internal/validators"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) CreateUser(user *models.User) error {
	// Validate the user input
	if err := validators.ValidateUser(user); err != nil {
		log.Printf("Validation failed: %v", err)
		return err
	}

	// Create a new user in the database
	if err := r.DB.Create(user).Error; err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	user := &models.User{}
	// Fetch user by ID
	if err := r.DB.First(user, id).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByVerificationToken(token string) (*models.User, error) {
	user := &models.User{}
	// Fetch user by verification token
	if err := r.DB.Where("verification_token = ?", token).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	// Update user in the database
	if err := r.DB.Save(user).Error; err != nil {
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
	if err := r.DB.Where("email = ?", email).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	// Fetch user by username
	if err := r.DB.Where("username = ?", username).First(user).Error; err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}
	return user, nil
}
