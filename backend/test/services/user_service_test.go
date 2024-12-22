package services

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/test/mocks"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockRepo *mocks.MockUserRepository
var service *services.UserService

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "true")
	code := m.Run()
	os.Exit(code)
}

func TestUserService_RegisterUser(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	// Define a user
	user := models.User{
		Username:           "joaquim",
		Email:              "example@gmail.com",
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

	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

	// Test the RegisterUser function
	err := service.RegisterUser(&user)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetUser tests the GetUser method of the UserService
func TestUserService_GetUser(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	expectedUser := models.User{
		Username:           "joaquim",
		Email:              "example@gmail.com",
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

	mockRepo.On("GetUserByID", uint(1)).Return(&expectedUser, nil)

	user, err := service.GetUser(1)

	assert.NoError(t, err)
	assert.Equal(t, &expectedUser, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetUserByInvalidID tests the GetUser method of the UserService with an invalid ID
func TestUserService_GetUserByInvalidID(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	mockRepo.On("GetUserByID", uint(1)).Return((*models.User)(nil), errors.New("user not found"))

	user, err := service.GetUser(1)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetUserByEmail tests the GetUserByEmail method of the UserService
func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	expectedUser := models.User{
		Username:           "joaquim",
		Email:              "example@gmail.com",
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

	mockRepo.On("GetUserByEmail", "example@gmail.com").Return(&expectedUser, nil)

	user, err := service.GetUserByEmail("example@gmail.com")

	assert.NoError(t, err)
	assert.Equal(t, &expectedUser, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetUserByInvalidEmail tests the GetUserByEmail method of the UserService with an invalid email
func TestUserService_GetUserByInvalidEmail(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	mockRepo.On("GetUserByEmail", "example@gmail.com").Return((*models.User)(nil), errors.New("user not found"))

	user, err := service.GetUserByEmail("example@gmail.com")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetUserByUsername tests the GetUserByUsername method of the UserService
func TestUserService_GetUserByUsername(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	expectedUser := models.User{
		Username:           "joaquim",
		Email:              "example@gmail.com",
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

	mockRepo.On("GetUserByUsername", "joaquim").Return(&expectedUser, nil)

	user, err := service.GetUserByUsername("joaquim")

	assert.NoError(t, err)
	assert.Equal(t, &expectedUser, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetUserByInvalidUsername tests the GetUserByUsername method of the UserService with an invalid username
func TestUserService_GetUserByInvalidUsername(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	mockRepo.On("GetUserByUsername", "joaquim").Return((*models.User)(nil), errors.New("user not found"))

	user, err := service.GetUserByUsername("joaquim")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_VerifyEmail tests the VerifyEmail method of the UserService
func TestUserService_VerifyEmail(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	user := models.User{
		Username:           "joaquim",
		Email:              "example@gmail.com",
		Password:           "123456",
		VerificationToken:  "validtoken",
		EmailVerified:      false,
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	user.CreatedAt = time.Now().Add(-1 * time.Minute)

	mockRepo.On("GetUserByVerificationToken", "validtoken").Return(&user, nil)
	mockRepo.On("UpdateUser", &user).Return(nil)

	err := service.VerifyEmail("validtoken")

	assert.NoError(t, err)
	assert.True(t, user.EmailVerified)
	mockRepo.AssertExpectations(t)
}

// TestUserService_VerifyEmailInvalidToken tests the VerifyEmail method of the UserService with an invalid token
func TestUserService_VerifyEmailInvalidToken(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	user := models.User{
		Username:           "joaquim",
		Email:              "example@gmail.com",
		Password:           "123456",
		VerificationToken:  "token",
		EmailVerified:      false,
		ProfilePicture:     "https://example.com/profile.jpg",
		Bio:                "I am a software engineer",
		Roles:              "admin,user",
		LastLogin:          time.Time{},
		DateOfBirth:        time.Time{},
		PreferredLanguage:  "en",
		ReadingPreferences: "novel,short_story",
		IsDeleted:          false,
	}

	user.CreatedAt = time.Now().Add(-20 * time.Minute)

	mockRepo.On("GetUserByVerificationToken", "token").Return(&user, nil)
	err := service.VerifyEmail("token")

	assert.Error(t, err)
	assert.Equal(t, "invalid token or token expired", err.Error())
	mockRepo.AssertNotCalled(t, "UpdateUser", mock.AnythingOfType("*models.User"))
	mockRepo.AssertExpectations(t)
}

// TestUserService_VerifyEmailInvalidUser tests the VerifyEmail method of the UserService with an invalid user
func TestUserService_VerifyEmailInvalidUser(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	mockRepo.On("GetUserByVerificationToken", "validtoken").Return((*models.User)(nil), errors.New("invalid user"))

	err := service.VerifyEmail("validtoken")

	assert.Error(t, err)
	assert.Equal(t, "invalid user", err.Error())
	mockRepo.AssertNotCalled(t, "UpdateUser", mock.AnythingOfType("*models.User"))
	mockRepo.AssertExpectations(t)
}

// TestUserService_DeleteUser tests the DeleteUser method of the UserService
func TestUserService_DeleteUser(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	user := models.User{
		Username:           "joaquim",
		Email:              "example@gmail.com",
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

	mockRepo.On("GetUserByID", uint(1)).Return(&user, nil)
	mockRepo.On("DeleteUser", &user).Return(nil)

	err := service.DeleteUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_DeleteInvalidUser tests the DeleteUser method of the UserService with an invalid user
func TestUserService_DeleteInvalidUser(t *testing.T) {
	mockRepo = new(mocks.MockUserRepository)
	service = services.NewUserService(mockRepo)

	mockRepo.On("GetUserByID", uint(1)).Return((*models.User)(nil), errors.New("user not found"))

	err := service.DeleteUser(1)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
