package controllers

import (
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestHandleRegister(t *testing.T) {
	tests := []struct {
		name                string
		description         string
		requestBody         string
		expectedCode        int
		expectedBody        string
		username            string
		email               string
		password            string
		bio                 string
		profilePicture      string
		dateOfBirth         string
		isAlreadyRegistered bool
	}{
		{
			name:                "#R_01",
			description:         "The user format is accepted",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_02",
			description:         "The user format is not accepted",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Invalid JSON"}`,
			requestBody:         `asdad`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_03",
			description:         "The user is already registered",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"User already registered"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: true,
		},
		{
			name:                "#R_04",
			description:         "There are more than one input",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Invalid JSON"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}, {"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_05",
			description:         "There are less than one input",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Invalid JSON"}`,
			requestBody:         ``,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_06",
			description:         "The input is not valid",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Invalid JSON"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_07",
			description:         "The username is more than 255 characters",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Username cannot be longer than 255 characters"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`, strings.Repeat("a", 256)),
			username:            strings.Repeat("a", 256),
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_08",
			description:         "The username is 255 characters",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`, strings.Repeat("a", 255)),
			username:            strings.Repeat("a", 255),
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_09",
			description:         "The username is empty",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Failed to bind JSON"}`,
			requestBody:         `{"username": "", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_10",
			description:         "The username is one character",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         `{"username": "j", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "j",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_11",
			description:         "The bio is more than 500 characters",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Bio cannot be longer than 500 characters"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "%s", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`, "joao", strings.Repeat("a", 501)),
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 strings.Repeat("a", 501),
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_12",
			description:         "The bio is 500 characters",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "%s", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`, "joao", strings.Repeat("a", 500)),
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 strings.Repeat("a", 500),
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_13",
			description:         "The bio is empty",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Failed to bind JSON"}`,
			requestBody:         `{"username": "joao", "bio": "", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_14",
			description:         "The bio is one character",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         `{"username": "joao", "bio": "b", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "b",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_15",
			description:         "The profile_picture is more than 255 characters",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Profile picture URL cannot be longer than 255 characters"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "bio", "profile_picture": "%s", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`, "joao", strings.Repeat("a", 256)),
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      strings.Repeat("a", 256),
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_16",
			description:         "The profile_picture is 255 characters",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "bio", "profile_picture": "%s", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`, "joao", strings.Repeat("a", 255)),
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      strings.Repeat("a", 255),
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_17",
			description:         "The profile_picture is empty",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Failed to bind JSON"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_18",
			description:         "The profile_picture is one character",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "p", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "p",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_19",
			description:         "The date_of_birth is less than 18 years",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"You must be at least 18 years old"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2008-01-01", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2008-01-01",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_20",
			description:         "The date_of_birth is 18 years",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2006-01-01", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2006-01-01",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_21",
			description:         "The date_of_birth is more than 18 years",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2000-01-01", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2000-01-01",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_22",
			description:         "The date_of_birth is in the future",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"You must be at least 18 years old"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2027-01-01", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2027-01-01",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_23",
			description:         "The email is more than 255 characters",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Email cannot be longer than 255 characters"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "%s@a.a", "password": "12345678"}`, "joao", strings.Repeat("a", 252)),
			username:            "joao",
			email:               fmt.Sprintf("%s@a.a", strings.Repeat("a", 252)),
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_24",
			description:         "The email is 255 characters",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "%s@gmail.com", "password": "12345678"}`, "joao", strings.Repeat("a", 245)),
			username:            "joao",
			email:               fmt.Sprintf("%s@gmail.com", strings.Repeat("a", 245)),
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_25",
			description:         "The email is empty",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Failed to bind JSON"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "", "password": "12345678"}`,
			username:            "joao",
			email:               "",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_26",
			description:         "The email is the smallest possible email",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "a@a.co", "password": "12345678"}`,
			username:            "joao",
			email:               "a@a.co",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_27",
			description:         "The password is less than 8 characters",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Password must be at least 8 characters long"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "1234567"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "1234567",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_28",
			description:         "The password is 8 characters",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         `{"username": "joao", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "12345678"}`,
			username:            "joao",
			email:               "joao@gmail.com",
			password:            "12345678",
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_29",
			description:         "The password is more than 72 characters",
			expectedCode:        http.StatusBadRequest,
			expectedBody:        `{"error":"Password cannot be longer than 72 characters"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "%s"}`, "joao", strings.Repeat("a", 73)),
			username:            "joao",
			email:               "joao@gmail.com",
			password:            strings.Repeat("a", 73),
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
		{
			name:                "#R_30",
			description:         "The password is 72 characters",
			expectedCode:        http.StatusCreated,
			expectedBody:        `{"message":"User registered successfully"}`,
			requestBody:         fmt.Sprintf(`{"username": "%s", "bio": "bio", "profile_picture": "profile_pic", "date_of_birth": "2004-12-23", "email": "joao@gmail.com", "password": "%s"}`, "joao", strings.Repeat("a", 72)),
			username:            "joao",
			email:               "joao@gmail.com",
			password:            strings.Repeat("a", 72),
			bio:                 "bio",
			profilePicture:      "profile_pic",
			dateOfBirth:         "2004-12-23",
			isAlreadyRegistered: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			_, err := userRepo.GetUserByEmail(tt.email)

			if err == nil {
				t.Errorf("Expected user to be deleted, but it was not")
			}

			// Mock user data
			user := models.User{
				Username:           tt.username,
				Email:              tt.email,
				Password:           tt.password,
				EmailVerified:      false,
				VerificationToken:  "",
				ProfilePicture:     tt.profilePicture,
				Bio:                tt.bio,
				Roles:              "admin,user",
				LastLogin:          time.Time{},
				DateOfBirth:        time.Time{},
				PreferredLanguage:  "en",
				ReadingPreferences: "novel,short_story",
				IsDeleted:          false,
			}

			if tt.isAlreadyRegistered {
				userRepo.CreateUser(&user)
			}

			// Create a request
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/auth/register", authController.Register)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Check if the database reflects changes
			if tt.expectedBody == `{"message":"User registered successfully"}` {
				createdUser, _ := userService.GetUser(1)

				if tt.isAlreadyRegistered {
					createdUser, _ = userService.GetUser(2)
				}
				assert.Equal(t, tt.username, createdUser.Username)
				assert.Equal(t, tt.email, createdUser.Email)
				assert.True(t, utils.ComparePassword(createdUser.Password, tt.password))
				assert.Equal(t, tt.bio, createdUser.Bio)
				assert.Equal(t, tt.profilePicture, createdUser.ProfilePicture)
				assert.Equal(t, tt.dateOfBirth, createdUser.DateOfBirth.Format("2006-01-02"))
			}
		})
	}
}

func TestHandleVerifyEmail(t *testing.T) {
	tests := []struct {
		name              string
		description       string
		requestBody       string
		expectedCode      int
		expectedBody      string
		queryParams       string
		verificationToken string
		isAlreadyVerified bool
		isAuthorized      bool
	}{
		{
			name:              "#V_01",
			description:       "The user is not found",
			expectedCode:      http.StatusNotFound,
			expectedBody:      `{"error":"User not found"}`,
			queryParams:       `token=123456789`,
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
		{
			name:              "#V_02",
			description:       "The user is created and not verified",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"Email verified successfully"}`,
			queryParams:       `token=123456789`,
			verificationToken: "123456789",
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
		{
			name:              "#V_03",
			description:       "The user is not authorized",
			expectedCode:      http.StatusUnauthorized,
			expectedBody:      `{"error":"Unauthorized"}`,
			queryParams:       `token=123456789`,
			verificationToken: "123456789",
			isAlreadyVerified: false,
			isAuthorized:      false,
		},
		{
			name:              "#V_04",
			description:       "The verification token is valid but the user is already verified",
			expectedCode:      http.StatusUnauthorized,
			expectedBody:      `{"error":"Invalid token or token expired"}`,
			queryParams:       `token=123456789`,
			verificationToken: "123456789",
			isAlreadyVerified: true,
			isAuthorized:      true,
		},
		{
			name:              "#V_05",
			description:       "There are more than one input, but has one valid",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"Email verified successfully"}`,
			queryParams:       `token=123456789&token2=123456789`,
			verificationToken: "123456789",
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
		{
			name:              "#V_06",
			description:       "There are less than one input",
			expectedCode:      http.StatusUnauthorized,
			expectedBody:      `{"error":"token is required"}`,
			verificationToken: "",
			queryParams:       ``,
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
		{
			name:              "#V_07",
			description:       "The input is not the right type",
			expectedCode:      http.StatusUnauthorized,
			expectedBody:      `{"error":"token is required"}`,
			queryParams:       `123`,
			verificationToken: "",
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
		{
			name:              "#V_08",
			description:       "The token is empty",
			expectedCode:      http.StatusUnauthorized,
			expectedBody:      `{"error":"token is required"}`,
			verificationToken: "",
			queryParams:       `token=`,
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
		{
			name:              "#V_09",
			description:       "The token is one character",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"Email verified successfully"}`,
			queryParams:       `token=a`,
			verificationToken: "a",
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
		{
			name:              "#V_10",
			description:       "The token is more than 255 characters",
			expectedCode:      http.StatusUnauthorized,
			expectedBody:      `{"error":"token cannot be longer than 255 characters"}`,
			queryParams:       fmt.Sprintf(`token=%s`, strings.Repeat("a", 256)),
			verificationToken: strings.Repeat("a", 256),
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
		{
			name:              "#V_11",
			description:       "The token is 255 characters",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"Email verified successfully"}`,
			queryParams:       fmt.Sprintf(`token=%s`, strings.Repeat("a", 255)),
			verificationToken: strings.Repeat("a", 255),
			isAlreadyVerified: false,
			isAuthorized:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			// Mock user data
			user := models.User{
				Username:           "joao",
				Email:              "joao@gmail.com",
				Password:           "12345678",
				EmailVerified:      tt.isAlreadyVerified,
				VerificationToken:  tt.verificationToken,
				ProfilePicture:     "profile_pic",
				Bio:                "bio",
				Roles:              "admin,user",
				LastLogin:          time.Time{},
				DateOfBirth:        time.Time{},
				PreferredLanguage:  "en",
				ReadingPreferences: "novel,short_story",
				IsDeleted:          false,
			}

			userRepo.CreateUser(&user)

			secretKey := os.Getenv("SECRET_KEY")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.Email,
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, _ := token.SignedString([]byte(secretKey))

			// Create a request
			req := httptest.NewRequest(http.MethodPost, "/auth/verify-email?"+tt.queryParams, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.isAuthorized {
				req.Header.Set("Authorization", "Bearer "+tokenString)
			}
			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/auth/verify-email", middleware.AuthMiddleware(), authController.VerifyEmail)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Check if the database reflects changes
			if tt.expectedBody == `{"message":"Email verified successfully"}` {
				validatedUser, _ := userService.GetUser(1)
				assert.True(t, validatedUser.EmailVerified)
				assert.Equal(t, validatedUser.VerificationToken, "")
			}
		})
	}
}

func TestHandleLogout(t *testing.T) {
	tests := []struct {
		name             string
		description      string
		requestBody      string
		expectedCode     int
		expectedBody     string
		refreshToken     string
		isAlreadyRevoked bool
		isAuthorized     bool
		isCreated        bool
	}{
		{
			name:         "#LO_01",
			description:  "The user is not found",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid token"}`,
			refreshToken: "valid_token",
			isCreated:    false,
			isAuthorized: true,
		},
		{
			name:         "#LO_02",
			description:  "The user is created and authorized",
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Logged out successfully"}`,
			refreshToken: "valid_token",
			isAuthorized: true,
			isCreated:    true,
		},
		{
			name:         "#LO_03",
			description:  "The user is not authorized",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Unauthorized"}`,
			refreshToken: "valid_token",
			isAuthorized: false,
			isCreated:    true,
		},
	}

	for _, tt := range tests {
		os.Setenv("SECRET_KEY", "your_secret_key")

		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			_, err := userRepo.GetUserByID(1)

			if err == nil {
				t.Errorf("Expected user to be deleted, but it was not")
			}

			// Mock user data
			pass, err := utils.HashPassword("12345678")

			if err != nil {
				t.Errorf("Failed to hash password: %v", err)
			}

			user := models.User{
				Username:           "joao",
				Email:              "joao@gmail.com",
				Password:           pass,
				EmailVerified:      true,
				VerificationToken:  "OLA",
				ProfilePicture:     "profile_pic",
				Bio:                "bio",
				Roles:              "admin,user",
				LastLogin:          time.Time{},
				DateOfBirth:        time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				PreferredLanguage:  "en",
				ReadingPreferences: "novel,short_story",
				IsDeleted:          false,
			}

			if tt.isCreated {
				userRepo.CreateUser(&user)
			}

			refreshToken := ""
			if tt.isCreated {
				userRepo.CreateUser(&user)
				createdUser, _ := userRepo.GetUserByID(1)
				_, refreshToken, _ = authService.GenerateToken(createdUser)
			}

			if tt.isAuthorized {
				tt.refreshToken = refreshToken
			}

			// Pass token in request header
			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.isAuthorized {
				req.Header.Set("Authorization", "Bearer "+tt.refreshToken)
			}

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/auth/logout", middleware.AuthMiddleware(), authController.Logout)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHandleLogin(t *testing.T) {
	tests := []struct {
		name              string
		description       string
		requestBody       string
		expectedCode      int
		expectedBody      string
		email             string
		password          string
		isAlreadyVerified bool
		isAuthorized      bool
		isFound           bool
	}{
		{
			name:              "#L_01",
			description:       "The user is not found",
			expectedCode:      http.StatusNotFound,
			expectedBody:      `{"error":"User not found"}`,
			email:             "joao@gmail.com",
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           false,
		},
		{
			name:              "#L_02",
			description:       "The user is created and not verified",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"User logged in successfully"}`,
			email:             "joao@gmail.com",
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_03",
			description:       "The user is not authorized",
			expectedCode:      http.StatusUnauthorized,
			expectedBody:      `{"error":"Unauthorized"}`,
			email:             "joao@gmail.com",
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      false,
			isFound:           true,
		},
		{
			name:              "#L_04",
			description:       "The user is authorized and the email is valid",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"User logged in successfully"}`,
			email:             "joao@gmail.com",
			password:          "12345678",
			isAlreadyVerified: true,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_05",
			description:       "There are more than one input, but has one valid",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"User logged in successfully"}`,
			email:             "joao@gmail.com",
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_06",
			description:       "There are less than one input",
			expectedCode:      http.StatusBadRequest,
			expectedBody:      `{"error":"Invalid input"}`,
			email:             "",
			password:          "",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_07",
			description:       "The email is empty",
			expectedCode:      http.StatusBadRequest,
			expectedBody:      `{"error":"Invalid input"}`,
			email:             "",
			password:          "",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_08",
			description:       "The email is smallest possible email",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"User logged in successfully"}`,
			email:             "j@a.co",
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_09",
			description:       "The email is more than 255 characters",
			expectedCode:      http.StatusBadRequest,
			expectedBody:      `{"error":"Email cannot be longer than 255 characters"}`,
			email:             fmt.Sprintf("%s@gmail.com", strings.Repeat("a", 246)),
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_10",
			description:       "The email is 255 characters",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"User logged in successfully"}`,
			email:             fmt.Sprintf("%s@gmail.com", strings.Repeat("a", 245)),
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_11",
			description:       "The email is one character",
			expectedCode:      http.StatusBadRequest,
			expectedBody:      `{"error":"Invalid email format"}`,
			email:             "a",
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_12",
			description:       "The password is less than 8 characters",
			expectedCode:      http.StatusBadRequest,
			expectedBody:      `{"error":"Password must be at least 8 characters long"}`,
			email:             "joao@gmail.com",
			password:          "1234567",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_13",
			description:       "The password is 8 characters",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"User logged in successfully"}`,
			email:             "joao@gmail.com",
			password:          "12345678",
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_14",
			description:       "The password is more than 72 characters",
			expectedCode:      http.StatusBadRequest,
			expectedBody:      `{"error":"Password cannot be longer than 72 characters"}`,
			email:             "joao@gmail.com",
			password:          strings.Repeat("a", 73),
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
		{
			name:              "#L_15",
			description:       "The password is 72 characters",
			expectedCode:      http.StatusOK,
			expectedBody:      `{"message":"User logged in successfully"}`,
			email:             "joao@gmail.com",
			password:          strings.Repeat("a", 72),
			isAlreadyVerified: false,
			isAuthorized:      true,
			isFound:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Setenv("SECRET_KEY", "your_secret_key")

			t.Cleanup(func() {
				cleanDB()
			})

			_, err := userRepo.GetUserByID(1)

			if err == nil {
				t.Errorf("Expected user to be deleted, but it was not")
			}

			// Mock user data
			pass, err := utils.HashPassword(tt.password)

			if err != nil || len(tt.password) > 72 || len(tt.password) < 8 {
				pass, _ = utils.HashPassword("12345678")
			}

			user := models.User{
				Username:           "joao",
				Email:              tt.email,
				Password:           pass,
				EmailVerified:      tt.isAlreadyVerified,
				VerificationToken:  "OLA",
				ProfilePicture:     "profile_pic",
				Bio:                "bio",
				Roles:              "admin,user",
				LastLogin:          time.Time{},
				DateOfBirth:        time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				PreferredLanguage:  "en",
				ReadingPreferences: "novel,short_story",
				IsDeleted:          false,
			}

			userRepo.CreateUser(&user)

			if !tt.isFound {
				tt.email = tt.email + "notfound"
			}

			secretKey := os.Getenv("SECRET_KEY")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.Email,
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, _ := token.SignedString([]byte(secretKey))

			request := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, tt.email, tt.password)

			// Create a request
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(request))
			req.Header.Set("Content-Type", "application/json")
			if tt.isAuthorized {
				req.Header.Set("Authorization", "Bearer "+tokenString)
			}

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/auth/login", middleware.AuthMiddleware(), authController.Login)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code != http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestHandleRefreshToken(t *testing.T) {
	tests := []struct {
		name             string
		description      string
		requestBody      string
		expectedCode     int
		expectedBody     string
		refreshToken     string
		isCreated        bool
		isAuthorized     bool
		isAlreadyRevoked bool
		isExpired        bool
		isFound          bool
		isLoggedIn       bool
	}{
		{
			name:         "#R_01",
			description:  "The user is not found",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid token"}`,
			refreshToken: "valid_token",
			isCreated:    false,
			isAuthorized: true,
			isFound:      true,
		},
		{
			name:         "#R_02",
			description:  "The user is created and not authorized",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Unauthorized"}`,
			refreshToken: "valid_token",
			isCreated:    true,
			isAuthorized: false,
			isFound:      true,
		},
		{
			name:         "#R_03",
			description:  "The user is not authorized",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Unauthorized"}`,
			refreshToken: "valid_token",
			isCreated:    true,
			isAuthorized: false,
			isFound:      true,
		},
		{
			name:             "#R_04",
			description:      "The refresh token is valid but is already revoked",
			expectedCode:     http.StatusUnauthorized,
			expectedBody:     `{"error":"Refresh token has been revoked"}`,
			refreshToken:     "valid_token",
			isCreated:        true,
			isAuthorized:     true,
			isFound:          true,
			isAlreadyRevoked: true,
		},
		{
			name:         "#R_05",
			description:  "The token is not expired",
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Logged in successfully"}`,
			refreshToken: "valid_token",
			isCreated:    true,
			isAuthorized: true,
			isFound:      true,
			isExpired:    false,
			isLoggedIn:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			os.Setenv("SECRET_KEY", "your_secret_key")

			_, err := userRepo.GetUserByID(1)

			if err == nil {
				t.Errorf("Expected user to be deleted, but it was not")
			}

			// Mock user data
			pass, _ := utils.HashPassword("12345678")

			user := models.User{
				Username:           "joao",
				Email:              "joao@gmail.com",
				Password:           pass,
				EmailVerified:      true,
				VerificationToken:  "OLA",
				ProfilePicture:     "profile_pic",
				Bio:                "bio",
				Roles:              "admin,user",
				LastLogin:          time.Time{},
				DateOfBirth:        time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				PreferredLanguage:  "en",
				ReadingPreferences: "novel,short_story",
				IsDeleted:          false,
			}

			userRepo.CreateUser(&user)

			refreshToken := ""
			if tt.isCreated {
				userRepo.CreateUser(&user)
				createdUser, _ := userRepo.GetUserByID(1)
				_, refreshToken, _ = authService.GenerateToken(createdUser)
			}

			if tt.isAuthorized {
				tt.refreshToken = refreshToken
			}

			if tt.isAlreadyRevoked {
				authRepo.RevokeToken(tt.refreshToken)
			}

			if tt.isLoggedIn {
				req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
				req.Header.Set("Content-Type", "application/json")
				if tt.isAuthorized {
					req.Header.Set("Authorization", "Bearer "+tt.refreshToken)
				}

				// Mock context and recorder
				w := httptest.NewRecorder()
				router := gin.Default()
				router.POST("/auth/login", middleware.AuthMiddleware(), authController.Login)
				router.ServeHTTP(w, req)
			}

			// Pass token in request header
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh-token", nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.isAuthorized {
				req.Header.Set("Authorization", "Bearer "+tt.refreshToken)
			}

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/auth/refresh-token", middleware.AuthMiddleware(), authController.RefreshToken)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code != http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
