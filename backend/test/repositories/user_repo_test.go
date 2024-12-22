package repositories

import (
	"backend/config"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/repositories/interfaces"
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB
var userRepo interfaces.UserRepositoryInterface

func cleanDB() {
	// Reset all tables by deleting rows (works for both SQLite and PostgreSQL)
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Failed to clean up the database: %v", err)
	}

	// For SQLite, reset auto-increment counters
	if db.Dialector.Name() == "sqlite" {
		if err := db.Exec("DELETE FROM sqlite_sequence WHERE name='users'").Error; err != nil {
			log.Fatalf("Failed to reset auto-increment for users table: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	// Setup: initialize the database connection and repositories
	db = config.ConnectDB(true) // true for test DB or in-memory DB
	userRepo = repositories.NewUserRepository(db)

	// Run tests
	code := m.Run()

	// Teardown: Cleanup any resources if necessary (e.g., closing DB connection)
	// For in-memory SQLite, nothing needs to be explicitly closed
	// If using a real DB, you might want to call db.Close() here

	// Exit with the code from running the tests
	os.Exit(code)
}

// TestUserRepository_CreateUser tests the CreateUser method of the UserRepository
func TestUserRepository_CreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	user := models.User{
		Username:           "joaquim",
		Email:              "joaquim@gmail.com",
		Password:           "123456",
		EmailVerified:      false,
		VerificationToken:  "123456",
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	if err := userRepo.CreateUser(&user); err != nil {
		t.Errorf("Failed to create user: %v", err)
	}
}

// TestUserRepository_GetUserByEmail tests the GetUserByEmail method of the UserRepository
func TestUserRepository_GetUserByEmail(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	user := models.User{
		Username:           "joaquim",
		Email:              "joaquim@gmail.com",
		Password:           "123456",
		EmailVerified:      false,
		VerificationToken:  "123456",
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	fetchedUser, err := userRepo.GetUserByEmail(user.Email)
	if err != nil {
		t.Errorf("Failed to fetch user: %v", err)
	}

	if fetchedUser.ID != user.ID {
		t.Errorf("Expected user ID to be %d, but got %d", user.ID, fetchedUser.ID)
	}
}

// TestUserRepository_GetUserByInvalidEmail tests the GetUserByEmail method of the UserRepository with an invalid email
func TestUserRepository_GetUserByInvalidEmail(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	fetchedUser, err := userRepo.GetUserByEmail("joel@gmail.com")
	if err == nil {
		t.Errorf("Expected error fetching user, but got nil")
	}
	if fetchedUser != nil {
		t.Errorf("Expected user to be nil, but got %v", fetchedUser)
	}
}

// TestUserRepository_GetUserByID tests the GetUserByID method of the UserRepository
func TestUserRepository_GetUserByID(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	user := models.User{
		Username:           "joaquim",
		Email:              "joaquim@gmail.com",
		Password:           "123456",
		EmailVerified:      false,
		VerificationToken:  "123456",
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	fetchedUser, err := userRepo.GetUserByID(uint(1))
	if err != nil {
		t.Errorf("Failed to fetch user: %v", err)
	}

	if fetchedUser.ID != 1 {
		t.Errorf("Expected user ID to be 1, but got %d", user.ID)
	}
}

// TestUserRepository_GetUserByInvalidID tests the GetUserByID method of the UserRepository with an invalid ID
func TestUserRepository_GetUserByInvalidID(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	fetchedUser, err := userRepo.GetUserByID(1)
	if err == nil {
		t.Errorf("Expected error fetching user, but got nil")
	}
	if fetchedUser != nil {
		t.Errorf("Expected user to be nil, but got %v", fetchedUser)
	}
}

// TestUserRepository_GetUserByUsername tests the GetUserByUsername method of the UserRepository
func TestUserRepository_GetUserByUsername(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	user := models.User{
		Username:           "joaquim",
		Email:              "joaquim@gmail.com",
		Password:           "123456",
		EmailVerified:      false,
		VerificationToken:  "123456",
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	fetchedUser, err := userRepo.GetUserByUsername(user.Username)
	if err != nil {
		t.Errorf("Failed to fetch user: %v", err)
	}

	if fetchedUser.ID != user.ID {
		t.Errorf("Expected user ID to be %d, but got %d", user.ID, fetchedUser.ID)
	}
}

// TestUserRepository_GetUserByInvalidUsername tests the GetUserByUsername method of the UserRepository with an invalid username
func TestUserRepository_GetUserByInvalidUsername(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	fetchedUser, err := userRepo.GetUserByUsername("joel")
	if err == nil {
		t.Errorf("Expected error fetching user, but got nil")
	}
	if fetchedUser != nil {
		t.Errorf("Expected user to be nil, but got %v", fetchedUser)
	}
}

// TestUserRepository_UpdateUser tests the UpdateUser method of the UserRepository
func TestUserRepository_UpdateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	user := models.User{
		Username:           "joaquim",
		Email:              "joaquim@gmail.com",
		Password:           "123456",
		EmailVerified:      false,
		VerificationToken:  "123456",
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	//Update User
	updatedUser := models.User{
		Username:           "joaquim2",
		Email:              "joaquim2@gmail.com",
		Password:           "1234567",
		EmailVerified:      true,
		VerificationToken:  "1234567",
		ProfilePicture:     "https://example.com/profile2.jpg",
		Bio:                "I am a software engineer2",
		Roles:              "admin",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "pt",
		ReadingPreferences: "novel,short_story,drama",
		IsDeleted:          true,
	}

	updatedUser.ID = 1

	// Update the user in the database
	if err := userRepo.UpdateUser(&updatedUser); err != nil {
		t.Errorf("Failed to update user: %v", err)
	}

	fetchedUser, err := userRepo.GetUserByID(uint(1))
	if err != nil {
		t.Errorf("Failed to fetch user: %v", err)
	}

	if fetchedUser.Username != updatedUser.Username {
		t.Errorf("Expected user username to be %s, but got %s", updatedUser.Username, fetchedUser.Username)
	}
	if fetchedUser.Email != updatedUser.Email {
		t.Errorf("Expected user email to be %s, but got %s", updatedUser.Email, fetchedUser.Email)
	}
	if fetchedUser.Password != updatedUser.Password {
		t.Errorf("Expected user password to be %s, but got %s", updatedUser.Password, fetchedUser.Password)
	}
	if fetchedUser.EmailVerified != updatedUser.EmailVerified {
		t.Errorf("Expected user emailVerified to be %t, but got %t", updatedUser.EmailVerified, fetchedUser.EmailVerified)
	}
	if fetchedUser.VerificationToken != updatedUser.VerificationToken {
		t.Errorf("Expected user verificationToken to be %s, but got %s", updatedUser.VerificationToken, fetchedUser.VerificationToken)
	}
	if fetchedUser.ProfilePicture != updatedUser.ProfilePicture {
		t.Errorf("Expected user profilePicture to be %s, but got %s", updatedUser.ProfilePicture, fetchedUser.ProfilePicture)
	}
	if fetchedUser.Bio != updatedUser.Bio {
		t.Errorf("Expected user bio to be %s, but got %s", updatedUser.Bio, fetchedUser.Bio)
	}
	if fetchedUser.Roles != updatedUser.Roles {
		t.Errorf("Expected user roles to be %s, but got %s", updatedUser.Roles, fetchedUser.Roles)
	}
	if fetchedUser.LastLogin != updatedUser.LastLogin {
		t.Errorf("Expected user lastLogin to be %s, but got %s", updatedUser.LastLogin, fetchedUser.LastLogin)
	}
	if fetchedUser.DateOfBirth != updatedUser.DateOfBirth {
		t.Errorf("Expected user dateOfBirth to be %s, but got %s", updatedUser.DateOfBirth, fetchedUser.DateOfBirth)
	}
	if fetchedUser.PreferredLanguage != updatedUser.PreferredLanguage {
		t.Errorf("Expected user preferredLanguage to be %s, but got %s", updatedUser.PreferredLanguage, fetchedUser.PreferredLanguage)
	}
	if fetchedUser.ReadingPreferences != updatedUser.ReadingPreferences {
		t.Errorf("Expected user readingPreferences to be %s, but got %s", updatedUser.ReadingPreferences, fetchedUser.ReadingPreferences)
	}
	if fetchedUser.IsDeleted != updatedUser.IsDeleted {
		t.Errorf("Expected user isDeleted to be %t, but got %t", updatedUser.IsDeleted, fetchedUser.IsDeleted)
	}
}

// TestUserRepository_UpdateInvalidUser tests the UpdateUser method of the UserRepository with an invalid user
func TestUserRepository_UpdateInvalidUser(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	//Update User
	updatedUser := models.User{
		Username:           "joaquim2",
		Email:              "joaquim2@gmail.com",
		Password:           "1234567",
		EmailVerified:      true,
		VerificationToken:  "1234567",
		ProfilePicture:     "https://example.com/profile2.jpg",
		Bio:                "I am a software engineer2",
		Roles:              "admin",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "pt",
		ReadingPreferences: "novel,short_story,drama",
		IsDeleted:          true,
	}

	updatedUser.ID = 1

	// Update the user in the database
	err := userRepo.UpdateUser(&updatedUser)

	if err != nil {
		t.Errorf("Expected an error updating the user, but got nil: %v", err)
	}
}

// TestUserRepository_DeleteUser tests the DeleteUser method of the UserRepository
func TestUserRepository_DeleteUser(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	user := models.User{
		Username:           "joaquim",
		Email:              "joaquim@gmail.com",
		Password:           "123456",
		EmailVerified:      false,
		VerificationToken:  "123456",
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := userRepo.DeleteUser(&user); err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}

	fetchedUser, err := userRepo.GetUserByID(uint(1))
	if err != nil {
		t.Errorf("Failed to fetch user: %v", err)
	}

	if fetchedUser.IsDeleted != true {
		t.Errorf("Expected user to be deleted, but it was not")
	}
}

// TestUserRepository_DeleteInvalidUser tests the DeleteUser method of the UserRepository with an invalid user
func TestUserRepository_DeleteInvalidUser(t *testing.T) {
	t.Cleanup(func() {
		cleanDB()
	})

	user := models.User{
		Username:           "joaquim",
		Email:              "joaquim@gmail.com",
		Password:           "123456",
		EmailVerified:      false,
		VerificationToken:  "123456",
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	err := userRepo.DeleteUser(&user)

	if err != nil {
		t.Errorf("Expected an error deleting the user, but got nil: %v", err)
	}
}
