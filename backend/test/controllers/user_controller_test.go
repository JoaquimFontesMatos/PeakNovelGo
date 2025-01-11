package controllers

import (
	"backend/config"
	"backend/internal/controllers"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/dtos"
	"backend/internal/repositories"
	repositoryInterfaces "backend/internal/repositories/interfaces"
	"backend/internal/services"
	serviceInterfaces "backend/internal/services/interfaces"
	"backend/internal/utils"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var db *gorm.DB
var userRepo repositoryInterfaces.UserRepositoryInterface
var authRepo repositoryInterfaces.AuthRepositoryInterface
var userService serviceInterfaces.UserServiceInterface
var authService serviceInterfaces.AuthServiceInterface
var userController controllers.UserController
var authController controllers.AuthController

func cleanDB() {
	os.Setenv("SECRET_KEY", "your_secret_key")

	// Reset all tables by deleting rows (works for both SQLite and PostgreSQL)
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Failed to clean up the database: %v", err)
	}

	if err := db.Exec("DELETE FROM revoked_tokens").Error; err != nil {
		log.Fatalf("Failed to clean up the database: %v", err)
	}

	// For SQLite, reset auto-increment counters
	if db.Dialector.Name() == "sqlite" {
		if err := db.Exec("DELETE FROM sqlite_sequence WHERE name='users'").Error; err != nil {
			log.Fatalf("Failed to reset auto-increment for users table: %v", err)
		}
		if err := db.Exec("DELETE FROM sqlite_sequence WHERE name='revoked_tokens'").Error; err != nil {
			log.Fatalf("Failed to reset auto-increment for revoked_tokens table: %v", err)
		}
	}

	userRepo = repositories.NewUserRepository(db)
	authRepo = repositories.NewAuthRepository(db)

	userService = services.NewUserService(userRepo)

	authService = services.NewAuthService(userRepo, authRepo)

	userController = *controllers.NewUserController(userService)

	authController = *controllers.NewAuthController(authService, userService)
}

func TestMain(m *testing.M) {
	// Setup: initialize the database connection and repositories
	db = config.ConnectDB(true) // true for test DB or in-memory DB

	cleanDB()

	// Setup Gin context
	gin.SetMode(gin.TestMode)

	os.Setenv("TESTING", "true")

	// Run tests
	code := m.Run()

	// Teardown: Cleanup any resources if necessary (e.g., closing DB connection)
	// For in-memory SQLite, nothing needs to be explicitly closed
	// If using a real DB, you might want to call db.Close() here

	// Exit with the code from running the tests
	os.Exit(code)
}

func TestHandleSoftDeleteUser(t *testing.T) {
	t.Run("#GSDU_01->The user is not created", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		// Call the function
		userController.HandleDeleteUser(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
	})

	t.Run("#GSDU_02->The user is created and the id is the correct one", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		// Call the function
		userController.HandleDeleteUser(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User deleted successfully")
	})

	t.Run("#GSDU_03->The user is not authorized", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		r := gin.Default()
		r.Use(middleware.AuthMiddleware()) // Attach the middleware
		r.GET("/users/id/:id", userController.HandleDeleteUser)

		// Create a test request with missing Authorization header
		req, _ := http.NewRequest("GET", "/users/id/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("#GSDU_04->The user is not found", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "2"}}

		// Call the function
		userController.HandleDeleteUser(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
	})

	t.Run("#GSDU_05->The user is found and authorized", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})
		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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
		userRepo.CreateUser(&user)

		// Generate a valid JWT token
		secretKey := os.Getenv("SECRET_KEY")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Email,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		validToken, _ := token.SignedString([]byte(secretKey))

		// Set up a test router with middleware
		r := gin.Default()
		r.Use(middleware.AuthMiddleware()) // Attach the middleware
		r.GET("/users/id/:id", userController.HandleDeleteUser)
		// Create a test request with missing Authorization header
		req, _ := http.NewRequest("GET", "/users/id/1", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User deleted successfully")
	})

	t.Run("#GSDU_07->There are less than one input", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// Call the function
		userController.HandleDeleteUser(c)
		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("#GSDU_08->user is already soft deleted", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})
		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
			EmailVerified:      false,
			VerificationToken:  "123456",
			ProfilePicture:     "https://example.com/profile.jpg",
			Bio:                "I am a software engineer",
			Roles:              "admin,user",
			LastLogin:          time.Time{},
			DateOfBirth:        time.Time{},
			PreferredLanguage:  "en",
			ReadingPreferences: "novel,short_story",
			IsDeleted:          true,
		}
		userRepo.CreateUser(&user)
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		// Call the function
		userController.HandleDeleteUser(c)
		// Assertions
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "User account is deactivated")
	})

	t.Run("#GSDU_09->Input is negative", func(t *testing.T) {
		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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
		userRepo.CreateUser(&user)
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "-1"}}
		// Call the function
		userController.HandleDeleteUser(c)
		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("#GSDU_10->The input is 0", func(t *testing.T) {
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "0"}}
		// Call the function
		userController.HandleDeleteUser(c)
		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("#GSDU_11->The input is a string", func(t *testing.T) {
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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
		userRepo.CreateUser(&user)
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "a"}}
		// Call the function
		userController.HandleDeleteUser(c)
		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("#GSDU_12->The input is max uint", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})
		user := models.User{
			Username:           "John Doe 2",
			Email:              "example2@mail.com",
			Password:           "12345678",
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
		var maxUint32 uint32 = math.MaxUint32
		user.ID = uint(maxUint32)
		err := userRepo.CreateUser(&user)
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}
		maxUint32Str := strconv.FormatUint(uint64(maxUint32), 10)
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: maxUint32Str}}
		// Call the function
		userController.HandleDeleteUser(c)
		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User deleted successfully")
	})

	t.Run("#GSDU_13->The input is more than max uint", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})
		user := models.User{
			Username:           "John Doe 2",
			Email:              "example2@mail.com",
			Password:           "12345678",
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
		var maxUint32 uint32 = math.MaxUint32
		maxUint32 += 1
		user.ID = uint(maxUint32)
		err := userRepo.CreateUser(&user)
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}
		maxUint32Str := strconv.FormatUint(uint64(maxUint32), 10)
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: maxUint32Str}}
		// Call the function
		userController.HandleDeleteUser(c)
		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})
}

func TestHandleGetUser(t *testing.T) {
	t.Run("#GUI_01->The user is not created", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch user")
	})
	t.Run("#GUI_02->The user is created and the id is the correct one", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_03->The user is not found", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch user")
	})

	t.Run("#GUI_04->There are more than one input", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}, {Key: "id2", Value: "2"}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Too many parameters")
	})

	t.Run("#GUI_05->There are less than one input", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "No parameters")
	})

	t.Run("#GUI_06->User is soft deleted", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
			EmailVerified:      false,
			VerificationToken:  "123456",
			ProfilePicture:     "https://example.com/profile.jpg",
			Bio:                "I am a software engineer",
			Roles:              "admin,user",
			LastLogin:          time.Time{},
			DateOfBirth:        time.Time{},
			PreferredLanguage:  "en",
			ReadingPreferences: "novel,short_story",
			IsDeleted:          true,
		}

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_07->Input is negative", func(t *testing.T) {
		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "-1"}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("#GUI_08->The input is 0", func(t *testing.T) {
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "0"}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("#GUI_09->The input is a string", func(t *testing.T) {
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "a"}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("#GUI_10->The input is max uint", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		user := models.User{
			Username:           "John Doe 2",
			Email:              "example2@mail.com",
			Password:           "12345678",
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

		user.ID = math.MaxUint32

		err := userRepo.CreateUser(&user)

		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}

		maxUint32 := math.MaxUint32

		maxUint32Str := strconv.FormatUint(uint64(maxUint32), 10)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: maxUint32Str}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe 2")
	})

	t.Run("#GUI_11->The input is more than max uint", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		user := models.User{
			Username:           "John Doe 2",
			Email:              "example2@mail.com",
			Password:           "12345678",
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

		var maxUint32 uint32 = math.MaxUint32
		maxUint32 += 1
		user.ID = uint(maxUint32)

		err := userRepo.CreateUser(&user)

		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}

		maxUint32Str := strconv.FormatUint(uint64(maxUint32), 10)
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: maxUint32Str}}

		// Call the function
		userController.HandleGetUser(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Run("#GUI_01->The user is not created", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: "example@mail.com"}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch user")
	})
	t.Run("#GUI_02->The user is created and the email is the correct one", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: "example@mail.com"}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_03->The user is not authorized", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		r := gin.Default()
		r.Use(middleware.AuthMiddleware()) // Attach the middleware
		r.GET("/users/email/:email", userController.HandleGetUserByEmail)

		// Create a test request with missing Authorization header
		req, _ := http.NewRequest("GET", "/users/email/example@mail.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("#GUI_04->The user is authorized and the email is valid", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock valid user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
			EmailVerified:      true,
			VerificationToken:  "123456",
			ProfilePicture:     "https://example.com/profile.jpg",
			Bio:                "I am a software engineer",
			Roles:              "admin,user",
			LastLogin:          time.Now(),
			DateOfBirth:        time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			PreferredLanguage:  "en",
			ReadingPreferences: "novel,short_story",
			IsDeleted:          false,
		}

		userRepo.CreateUser(&user)

		// Generate a valid JWT token
		secretKey := os.Getenv("SECRET_KEY")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Email,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		validToken, _ := token.SignedString([]byte(secretKey))

		// Set up a test router with middleware
		r := gin.Default()
		r.Use(middleware.AuthMiddleware()) // Attach the middleware
		r.GET("/users/email/:email", userController.HandleGetUserByEmail)

		// Create a test request with a valid Authorization header
		req, _ := http.NewRequest("GET", "/users/email/example@mail.com", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_05->The user is not found", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: "example@mail.com"}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch user")
	})

	t.Run("#GUI_06->There are more than one input", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: "example@mail.com"}, {Key: "email2", Value: "example2@mail.com"}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Too many parameters")
	})

	t.Run("#GUI_07->There are less than one input", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "No parameters")
	})

	t.Run("#GUI_08->User is soft deleted", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
			EmailVerified:      false,
			VerificationToken:  "123456",
			ProfilePicture:     "https://example.com/profile.jpg",
			Bio:                "I am a software engineer",
			Roles:              "admin,user",
			LastLogin:          time.Time{},
			DateOfBirth:        time.Time{},
			PreferredLanguage:  "en",
			ReadingPreferences: "novel,short_story",
			IsDeleted:          true,
		}

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: "example@mail.com"}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_09->Email is not valid", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: "examplemail.com"}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email format")
	})

	t.Run("#GUI_10->Email is empty", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: ""}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Email is required")
	})

	t.Run("#GUI_11->Email is 255 characters", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasddasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda@gmail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: "asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasddasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda@gmail.com"}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_12->Email is more than 255 characters", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "dasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasddasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda@gmail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "email", Value: "dasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasddasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda@gmail.com"}}

		// Call the function
		userController.HandleGetUserByEmail(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Email cannot be longer than 255 characters")
	})
}

func TestGetUserByUsername(t *testing.T) {
	t.Run("#GUI_01->The user is not created", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "username", Value: "John Doe"}}

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch user")
	})
	t.Run("#GUI_02->The user is created and the username is the correct one", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "username", Value: "John Doe"}}

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_03->The user is not authorized", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		r := gin.Default()
		r.Use(middleware.AuthMiddleware()) // Attach the middleware
		r.GET("/users/username/:username", userController.HandleGetUserByUsername)

		// Create a test request with missing Authorization header
		req, _ := http.NewRequest("GET", "/users/username/John Doe", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("#GUI_04->The user is not found", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "username", Value: "John Doe"}}

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch user")
	})

	t.Run("#GUI_05->The user is found and authorized", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Generate a valid JWT token
		secretKey := os.Getenv("SECRET_KEY")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Email,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		validToken, _ := token.SignedString([]byte(secretKey))

		// Set up a test router with middleware
		r := gin.Default()
		r.Use(middleware.AuthMiddleware()) // Attach the middleware
		r.GET("/users/username/:username", userController.HandleGetUserByUsername)

		// Create a test request with a valid Authorization header
		req, _ := http.NewRequest("GET", "/users/username/John Doe", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_07->There are less than one input", func(t *testing.T) {
		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Username is required")
	})

	t.Run("#GUI_08->User is soft deleted", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
			EmailVerified:      false,
			VerificationToken:  "123456",
			ProfilePicture:     "https://example.com/profile.jpg",
			Bio:                "I am a software engineer",
			Roles:              "admin,user",
			LastLogin:          time.Time{},
			DateOfBirth:        time.Time{},
			PreferredLanguage:  "en",
			ReadingPreferences: "novel,short_story",
			IsDeleted:          true,
		}

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "username", Value: "John Doe"}}

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("#GUI_09->Username is empty", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "John Doe",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "username", Value: ""}}

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Username is required")
	})

	t.Run("#GUI_10->Username is one character", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "a",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "username", Value: "a"}}

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"username\":\"a\"")
	})

	t.Run("#GUI_11->Username is more than 255 characters", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasd",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "username", Value: "dasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasd"}}

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Username cannot be longer than 255 characters")
	})

	t.Run("#GUI_12->Username is 255 characters", func(t *testing.T) {
		t.Cleanup(func() {
			cleanDB()
		})

		// Mock user data
		user := models.User{
			Username:           "asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasd",
			Email:              "example@mail.com",
			Password:           "12345678",
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

		userRepo.CreateUser(&user)

		// Create a test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "username", Value: "asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasd"}}

		// Call the function
		userController.HandleGetUserByUsername(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasd")
	})
}

func TestHandleEmailUpdate(t *testing.T) {
	tests := []struct {
		name         string
		description  string
		urlParam     string
		requestBody  string
		expectedCode int
		expectedBody string
		createUser   bool
		isAuthorized bool
		newEmail     string
		IsDeleted    bool
	}{
		{
			name:         "#EUE_01",
			description:  "User not created",
			urlParam:     "1",
			requestBody:  `{"new_email": "new@example.com"}`,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"User not found"}`,
			createUser:   false,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "new@example.com",
		},
		{
			name:         "#EUE_02",
			description:  "The user is created and the id is valid",
			urlParam:     "1",
			requestBody:  `{"new_email": "new@example.com"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Email update initiated. Please verify the new email."}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "new@example.com",
		},
		{
			name:         "#EUE_03",
			description:  "The user is not authorized",
			urlParam:     "1",
			requestBody:  `{"new_email": "new@example.com"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Unauthorized"}`,
			createUser:   true,
			isAuthorized: false,
			IsDeleted:    false,
			newEmail:     "new@example.com",
		},
		{
			name:         "#EUE_04",
			description:  "The user is authorized and the email is valid",
			urlParam:     "1",
			requestBody:  `{"new_email": "new@example.com"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Email update initiated. Please verify the new email."}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "new@example.com",
		},
		{
			name:         "#EUE_05",
			description:  "The user is not found",
			urlParam:     "2",
			requestBody:  `{"new_email": "new@example.com"}`,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"User not found"}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "new@example.com",
		},
		{
			name:         "#EUE_06",
			description:  "There are more than one id",
			urlParam:     "1,4",
			requestBody:  `{"new_email": "new@example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Invalid ID"}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "new@example.com",
		},
		{
			name:         "#EUE_07",
			description:  "There are less than one id",
			requestBody:  `{"new_email": "new@example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Invalid ID"}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "",
		},
		{
			name:         "#EUE_08",
			description:  "User is soft deleted",
			urlParam:     "1",
			requestBody:  `{"new_email": "new@example.com"}`,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"User account is deactivated"}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    true,
			newEmail:     "new@example.com",
		},
		{
			name:         "#EUE_09",
			description:  "Email is not valid",
			urlParam:     "1",
			requestBody:  `{"new_email": "newxample.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Failed to bind JSON"}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "newxample.com",
		},
		{
			name:         "#EUE_10",
			description:  "Email is empty",
			urlParam:     "1",
			requestBody:  `{"new_email": ""}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Failed to bind JSON"}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "",
		},
		{
			name:         "#EUE_11",
			description:  "Email is 255 characters",
			urlParam:     "1",
			requestBody:  `{"new_email": "asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasddasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda@gmail.com"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Email update initiated. Please verify the new email."}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasddasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda@gmail.com",
		},
		{
			name:         "#EUE_12",
			description:  "Email is more than 255 characters",
			urlParam:     "1",
			requestBody:  `{"new_email": "dasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasddasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda@gmail.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Email cannot be longer than 255 characters"}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "dasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasddasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda@gmail.com",
		},
		{
			name:         "#EUE_13",
			description:  "The input is the wrong type",
			urlParam:     "1",
			requestBody:  `{"new_email": 1}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Failed to bind JSON"}`,
			createUser:   true,
			isAuthorized: true,
			IsDeleted:    false,
			newEmail:     "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			_, err := userRepo.GetUserByID(1)

			if err == nil {
				t.Errorf("Expected user to be deleted, but it was not")
			}

			// Mock user data
			user := models.User{
				Username:           "example",
				Email:              "test@example.com",
				Password:           "12345678",
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

			if tt.IsDeleted {
				user.IsDeleted = true
			}

			if tt.createUser {
				userRepo.CreateUser(&user)
			}

			secretKey := os.Getenv("SECRET_KEY")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.Email,
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, _ := token.SignedString([]byte(secretKey))

			// Create a request
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.urlParam+"/email", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.isAuthorized {
				req.Header.Set("Authorization", "Bearer "+tokenString)
			}

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.PUT("/users/:id/email", middleware.AuthMiddleware(), userController.UpdateEmail)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Check if the database reflects changes
			if tt.expectedBody == `{"message":"Email update initiated. Please verify the new email."}` {
				updatedUser, _ := userService.GetUser(1)
				assert.Equal(t, tt.newEmail, updatedUser.Email)
				assert.NotEqual(t, user.VerificationToken, updatedUser.VerificationToken)
			}
		})
	}
}

func TestHandleUpdatePassword(t *testing.T) {
	tests := []struct {
		name         string
		description  string
		urlParam     string
		requestBody  string
		expectedCode int
		expectedBody string
		createUser   bool
		isAuthorized bool
		newPassword  string
		isDeleted    bool
	}{
		{
			name:         "#EUP_01",
			description:  "User not created",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassword"}`,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"User not found"}`,
			createUser:   false,
			isAuthorized: true,
			newPassword:  "newpassword",
			isDeleted:    false,
		},
		{
			name:         "#EUP_02",
			description:  "The user is created and the id is valid",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassword"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Password updated successfully"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "newpassword",
			isDeleted:    false,
		},
		{
			name:         "#EUP_03",
			description:  "The user is not authorized",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassword"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Unauthorized"}`,
			createUser:   true,
			isAuthorized: false,
			newPassword:  "newpassword",
			isDeleted:    false,
		},
		{
			name:         "#EUP_04",
			description:  "The user is authorized and the email is valid",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassword"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Password updated successfully"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "newpassword",
			isDeleted:    false,
		},
		{
			name:         "#EUP_05",
			description:  "The user is not found",
			urlParam:     "2",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassword"}`,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"User not found"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "newpassword",
			isDeleted:    false,
		},
		{
			name:         "#EUP_06",
			description:  "There are more than one id",
			urlParam:     "1,4",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassword"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Invalid ID"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "newpassword",
			isDeleted:    false,
		},
		{
			name:         "#EUP_07",
			description:  "There are less than one id",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassword"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Invalid ID"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "newpassword",
			isDeleted:    false,
		},
		{
			name:         "#EUP_08",
			description:  "User is soft deleted",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassword"}`,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"User account is deactivated"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "newpassword",
			isDeleted:    true,
		},
		{
			name:         "#EUP_09",
			description:  "Current password is incorrect",
			urlParam:     "1",
			requestBody:  `{"current_password": "123456789", "new_password": "newpassword"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid password"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "newpassword",
			isDeleted:    false,
		},
		{
			name:         "#EUP_10",
			description:  "New password is the same as the current password",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": "12345678"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"New password cannot be the same as the current password"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "12345678",
			isDeleted:    false,
		},
		{
			name:         "#EUP_11",
			description:  "New password is not valid (too short)",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": "1234567"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Password must be at least 8 characters long"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "1234567",
			isDeleted:    false,
		},
		{
			name:         "#EUP_12",
			description:  "New password is 8 characters",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": "newpassw"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Password updated successfully"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "newpassw",
			isDeleted:    false,
		},
		{
			name:         "#EUP_13",
			description:  "New password is 72 characters",
			urlParam:     "1",
			requestBody:  fmt.Sprintf(`{"current_password": "12345678", "new_password": "%s"}`, strings.Repeat("a", 72)),
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Password updated successfully"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  strings.Repeat("a", 72),
			isDeleted:    false,
		},
		{
			name:         "#EUP_14",
			description:  "New password is more than 72 characters (73)",
			urlParam:     "1",
			requestBody:  fmt.Sprintf(`{"current_password": "12345678", "new_password": "%s"}`, strings.Repeat("a", 73)),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Password cannot be longer than 72 characters"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  strings.Repeat("a", 73),
			isDeleted:    false,
		},
		{
			name:         "#EUP_15",
			description:  "New password is empty",
			urlParam:     "1",
			requestBody:  `{"current_password": "12345678", "new_password": ""}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Failed to bind JSON"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "",
			isDeleted:    false,
		},
		{
			name:         "#EUP_16",
			description:  "Input is invalid",
			urlParam:     "1",
			requestBody:  `a`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Invalid JSON"}`,
			createUser:   true,
			isAuthorized: true,
			newPassword:  "1",
			isDeleted:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			_, err := userRepo.GetUserByID(1)

			if err == nil {
				t.Errorf("Expected user to be deleted, but it was not")
			}

			// Mock user data
			user := models.User{
				Username:    "example",
				Email:       "test@example.com",
				Password:    "12345678",
				Bio:         "I am a software engineer",
				LastLogin:   time.Time{},
				DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			}

			userFields := dtos.RegisterRequest{
				Username:       user.Username,
				Email:          user.Email,
				Password:       user.Password,
				Bio:            user.Bio,
				ProfilePicture: user.ProfilePicture,
				DateOfBirth:    user.DateOfBirth.Format("2006-01-02"),
			}

			if tt.createUser {
				userService.RegisterUser(&userFields)
			}

			if tt.isDeleted {
				userService.DeleteUser(1)
			}

			secretKey := os.Getenv("SECRET_KEY")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.Email,
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, _ := token.SignedString([]byte(secretKey))

			// Create a request
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.urlParam+"/password", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.isAuthorized {
				req.Header.Set("Authorization", "Bearer "+tokenString)
			}

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.PUT("/users/:id/password", middleware.AuthMiddleware(), userController.UpdatePassword)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			// Check if the database reflects changes
			if tt.expectedBody == `{"message":"Password updated successfully"}` {
				updatedUser, _ := userService.GetUser(1)
				assert.True(t, utils.ComparePassword(updatedUser.Password, tt.newPassword))
			}
		})
	}
}

func TestHandleUpdateFields(t *testing.T) {
	tests := []struct {
		name                  string
		description           string
		urlParam              string
		requestBody           string
		expectedCode          int
		expectedBody          string
		createUser            bool
		isAuthorized          bool
		newUsername           string
		newBio                string
		newProfilePicture     string
		newPreferredLanguage  string
		newReadingPreferences string
		newDateOfBirth        string
		newRoles              string
		IsDeleted             bool
	}{
		{
			name:                  "#EUF_01",
			description:           "User not created",
			urlParam:              "1",
			requestBody:           `{"username": "newusername", "bio": "newbio", "profile_picture": "https://example.com/profile.jpg", "preferred_language": "en", "reading_preferences": "novel,short_story", "date_of_birth": "1990-01-01", "roles": "admin,user"}`,
			expectedCode:          http.StatusNotFound,
			expectedBody:          `{"error":"User not found"}`,
			createUser:            false,
			isAuthorized:          true,
			newUsername:           "newusername",
			newBio:                "newbio",
			newProfilePicture:     "https://example.com/profile.jpg",
			newPreferredLanguage:  "en",
			newReadingPreferences: "novel,short_story",
			newDateOfBirth:        "1990-01-01",
			newRoles:              "admin,user",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_02",
			description:           "The user is created and validated",
			urlParam:              "1",
			requestBody:           `{"username": "newusername", "bio": "newbio", "profile_picture": "https://example.com/profile.jpg", "preferred_language": "en", "reading_preferences": "novel,short_story", "date_of_birth": "1990-01-01", "roles": "admin,user"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "newusername",
			newBio:                "newbio",
			newProfilePicture:     "https://example.com/profile.jpg",
			newPreferredLanguage:  "en",
			newReadingPreferences: "novel,short_story",
			newDateOfBirth:        "1990-01-01",
			newRoles:              "admin,user",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_03",
			description:           "The user is not authorized",
			urlParam:              "1",
			requestBody:           `{"username": "newusername", "bio": "newbio", "profile_picture": "https://example.com/profile.jpg", "preferred_language": "en", "reading_preferences": "novel,short_story", "date_of_birth": "1990-01-01", "roles": "admin,user"}`,
			expectedCode:          http.StatusUnauthorized,
			expectedBody:          `{"error":"Unauthorized"}`,
			createUser:            true,
			isAuthorized:          false,
			newUsername:           "newusername",
			newBio:                "newbio",
			newProfilePicture:     "https://example.com/profile.jpg",
			newPreferredLanguage:  "en",
			newReadingPreferences: "novel,short_story",
			newDateOfBirth:        "1990-01-01",
			newRoles:              "admin,user",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_04",
			description:           "The user is not found",
			urlParam:              "2",
			requestBody:           `{"username": "newusername", "bio": "newbio", "profile_picture": "https://example.com/profile.jpg", "preferred_language": "en", "reading_preferences": "novel,short_story", "date_of_birth": "1990-01-01", "roles": "admin,user"}`,
			expectedCode:          http.StatusNotFound,
			expectedBody:          `{"error":"User not found"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "newusername",
			newBio:                "newbio",
			newProfilePicture:     "https://example.com/profile.jpg",
			newPreferredLanguage:  "en",
			newReadingPreferences: "novel,short_story",
			newDateOfBirth:        "1990-01-01",
			newRoles:              "admin,user",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_05",
			description:           "There are more than one input",
			urlParam:              "1,4",
			requestBody:           `{"username": "newusername", "bio": "newbio", "profile_picture": "https://example.com/profile.jpg", "preferred_language": "en", "reading_preferences": "novel,short_story", "date_of_birth": "1990-01-01", "roles": "admin,user"}`,
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Invalid ID"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "newusername",
			newBio:                "newbio",
			newProfilePicture:     "https://example.com/profile.jpg",
			newPreferredLanguage:  "en",
			newReadingPreferences: "novel,short_story",
			newDateOfBirth:        "1990-01-01",
			newRoles:              "admin,user",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_06",
			description:           "There are less than one input",
			requestBody:           `{"username": "newusername", "bio": "newbio", "profile_picture": "https://example.com/profile.jpg", "preferred_language": "en", "reading_preferences": "novel,short_story", "date_of_birth": "1990-01-01", "roles": "admin,user"}`,
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Invalid ID"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "newusername",
			newBio:                "newbio",
			newProfilePicture:     "https://example.com/profile.jpg",
			newPreferredLanguage:  "en",
			newReadingPreferences: "novel,short_story",
			newDateOfBirth:        "1990-01-01",
			newRoles:              "admin,user",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_07",
			description:           "User is soft deleted",
			urlParam:              "1",
			requestBody:           `{"username": "newusername", "bio": "newbio", "profile_picture": "https://example.com/profile.jpg", "preferred_language": "en", "reading_preferences": "novel,short_story", "date_of_birth": "1990-01-01", "roles": "admin,user"}`,
			expectedCode:          http.StatusForbidden,
			expectedBody:          `{"error":"User account is deactivated"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "newusername",
			newBio:                "newbio",
			newProfilePicture:     "https://example.com/profile.jpg",
			newPreferredLanguage:  "en",
			newReadingPreferences: "novel,short_story",
			newDateOfBirth:        "1990-01-01",
			newRoles:              "admin,user",
			IsDeleted:             true,
		},
		{
			name:                  "#EUF_08",
			description:           "Fields is not valid",
			urlParam:              "1",
			requestBody:           `a`,
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Invalid JSON"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_09",
			description:           "Fields is empty",
			urlParam:              "1",
			requestBody:           `{}`,
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"No fields provided"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_10",
			description:           "Not every field filled",
			urlParam:              "1",
			requestBody:           `{"username": "newusername"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "newusername",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_11",
			description:           "Fields username is more than 255 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"username": "%s"}`, strings.Repeat("a", 256)),
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Username cannot be longer than 255 characters"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           strings.Repeat("a", 256),
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_12",
			description:           "Fields username is 255 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"username": "%s"}`, strings.Repeat("a", 255)),
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           strings.Repeat("a", 255),
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_13",
			description:           "Fields username is one character",
			urlParam:              "1",
			requestBody:           `{"username": "a"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "a",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_14",
			description:           "Fields bio is more than 500 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"bio": "%s"}`, strings.Repeat("a", 501)),
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Bio cannot be longer than 500 characters"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                strings.Repeat("a", 501),
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_15",
			description:           "Fields bio is 500 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"bio": "%s"}`, strings.Repeat("a", 500)),
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                strings.Repeat("a", 500),
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_16",
			description:           "Fields bio is one character",
			urlParam:              "1",
			requestBody:           `{"bio": "a"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "a",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_17",
			description:           "Fields profile_picture is more than 255 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"profile_picture": "%s"}`, strings.Repeat("a", 256)),
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Profile picture URL cannot be longer than 255 characters"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     strings.Repeat("a", 256),
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_18",
			description:           "Fields profile_picture is 255 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"profile_picture": "%s"}`, strings.Repeat("a", 255)),
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     strings.Repeat("a", 255),
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_19",
			description:           "Fields profile_picture is one character",
			urlParam:              "1",
			requestBody:           `{"profile_picture": "a"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "a",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_20",
			description:           "Fields preferred_language is more than 100 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"preferred_language": "%s"}`, strings.Repeat("a", 101)),
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Preferred language cannot be longer than 100 characters"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  strings.Repeat("a", 101),
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_21",
			description:           "Fields preferred_language is 100 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"preferred_language": "%s"}`, strings.Repeat("a", 100)),
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  strings.Repeat("a", 100),
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_22",
			description:           "Fields preferred_language is one character",
			urlParam:              "1",
			requestBody:           `{"preferred_language": "a"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "a",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_23",
			description:           "Fields reading_preferences is more than 255 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"reading_preferences": "%s"}`, strings.Repeat("a", 256)),
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Reading preferences cannot be longer than 255 characters"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: strings.Repeat("a", 256),
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_24",
			description:           "Fields reading_preferences is 255 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"reading_preferences": "%s"}`, strings.Repeat("a", 255)),
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: strings.Repeat("a", 255),
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_25",
			description:           "Fields reading_preferences is one character",
			urlParam:              "1",
			requestBody:           `{"reading_preferences": "a"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "a",
			newDateOfBirth:        "",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_26",
			description:           "Fields date_of_birth is less than 18 years",
			urlParam:              "1",
			requestBody:           `{"date_of_birth": "2009-01-01"}`,
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"You must be at least 18 years old"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "2009-01-01",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_27",
			description:           "Fields date_of_birth is 18 years",
			urlParam:              "1",
			requestBody:           `{"date_of_birth": "2006-01-01"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "2006-01-01",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_28",
			description:           "Fields date_of_birth is more than 18 years",
			urlParam:              "1",
			requestBody:           `{"date_of_birth": "2000-01-01"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "2000-01-01",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_29",
			description:           "Fields date_of_birth is in the future",
			urlParam:              "1",
			requestBody:           `{"date_of_birth": "2027-01-01"}`,
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"You must be at least 18 years old"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "2027-01-01",
			newRoles:              "",
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_30",
			description:           "Fields roles is more than 255 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"roles": "%s"}`, strings.Repeat("a", 256)),
			expectedCode:          http.StatusBadRequest,
			expectedBody:          `{"error":"Roles cannot be longer than 255 characters"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              strings.Repeat("a", 256),
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_31",
			description:           "Fields roles is 255 characters",
			urlParam:              "1",
			requestBody:           fmt.Sprintf(`{"roles": "%s"}`, strings.Repeat("a", 255)),
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              strings.Repeat("a", 255),
			IsDeleted:             false,
		},
		{
			name:                  "#EUF_32",
			description:           "Fields roles is one character",
			urlParam:              "1",
			requestBody:           `{"roles": "a"}`,
			expectedCode:          http.StatusOK,
			expectedBody:          `{"message":"User updated successfully"}`,
			createUser:            true,
			isAuthorized:          true,
			newUsername:           "",
			newBio:                "",
			newProfilePicture:     "",
			newPreferredLanguage:  "",
			newReadingPreferences: "",
			newDateOfBirth:        "",
			newRoles:              "a",
			IsDeleted:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			_, err := userRepo.GetUserByID(1)

			if err == nil {
				t.Errorf("Expected user to be deleted, but it was not")
			}

			// Mock user data
			user := models.User{
				Username:           "example",
				Email:              "test@example.com",
				Password:           "12345678",
				EmailVerified:      false,
				VerificationToken:  "",
				ProfilePicture:     "https://example.com/profile2.jpg",
				Bio:                "I am a software engineer",
				Roles:              "admin2,user",
				LastLogin:          time.Time{},
				DateOfBirth:        time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
				PreferredLanguage:  "ru",
				ReadingPreferences: "novel,short_story,movis",
				IsDeleted:          false,
			}

			if tt.IsDeleted {
				user.IsDeleted = true
			}

			if tt.createUser {
				userRepo.CreateUser(&user)
			}

			secretKey := os.Getenv("SECRET_KEY")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.Email,
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, _ := token.SignedString([]byte(secretKey))

			// Create a request
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.urlParam+"/fields", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.isAuthorized {
				req.Header.Set("Authorization", "Bearer "+tokenString)
			}

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.PUT("/users/:id/fields", middleware.AuthMiddleware(), userController.UpdateUserFields)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Check if the database reflects changes
			if tt.expectedBody == `{"message":"User updated successfully"}` {
				updatedUser, _ := userService.GetUser(1)
				if tt.newUsername != "" {
					assert.Equal(t, tt.newUsername, updatedUser.Username)
				} else {
					assert.Equal(t, user.Username, updatedUser.Username)
				}
				if tt.newBio != "" {
					assert.Equal(t, tt.newBio, updatedUser.Bio)
				} else {
					assert.Equal(t, user.Bio, updatedUser.Bio)
				}
				if tt.newProfilePicture != "" {
					assert.Equal(t, tt.newProfilePicture, updatedUser.ProfilePicture)
				} else {
					assert.Equal(t, user.ProfilePicture, updatedUser.ProfilePicture)
				}
				if tt.newPreferredLanguage != "" {
					assert.Equal(t, tt.newPreferredLanguage, updatedUser.PreferredLanguage)
				} else {
					assert.Equal(t, user.PreferredLanguage, updatedUser.PreferredLanguage)
				}
				if tt.newReadingPreferences != "" {
					assert.Equal(t, tt.newReadingPreferences, updatedUser.ReadingPreferences)
				} else {
					assert.Equal(t, user.ReadingPreferences, updatedUser.ReadingPreferences)
				}
				if tt.newDateOfBirth != "" {
					assert.Equal(t, tt.newDateOfBirth, updatedUser.DateOfBirth.Format("2006-01-02"))
				} else {
					assert.Equal(t, user.DateOfBirth.Format("2006-01-02"), updatedUser.DateOfBirth.Format("2006-01-02"))
				}
				if tt.newRoles != "" {
					assert.Equal(t, tt.newRoles, updatedUser.Roles)
				} else {
					assert.Equal(t, user.Roles, updatedUser.Roles)
				}
			}
		})
	}
}
