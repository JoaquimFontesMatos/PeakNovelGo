package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"backend/internal/dtos"
	"backend/internal/services/interfaces"
	"backend/internal/types"
	"backend/internal/validators"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type AuthController struct {
	AuthService interfaces.AuthServiceInterface
	UserService interfaces.UserServiceInterface
}

func NewAuthController(authService interfaces.AuthServiceInterface, userService interfaces.UserServiceInterface) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserService: userService,
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var registerRequest dtos.RegisterRequest

	if err := validators.ValidateBody(c, &registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the user service to create a new user
	if err := ac.UserService.RegisterUser(&registerRequest); err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			case types.CONFLICT_ERROR:
				c.JSON(http.StatusConflict, gin.H{"error": myError.Message})
			}
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req dtos.LoginRequest // Define a login request struct

	// Bind the incoming JSON to the LoginRequest model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate credentials
	user, err := ac.AuthService.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.INTERNAL_SERVER_ERROR:
				c.JSON(http.StatusInternalServerError, gin.H{"error": myError.Message})
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			}
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT tokens (access token and refresh token)
	accessToken, refreshToken, err := ac.AuthService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userDto, err := dtos.ConvertUserModelToDTO(*user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the access token and user info
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         userDto,
	})
}

func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid refresh token header"})
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Call the service to refresh the token
	newAccessToken, newRefreshToken, user, err := ac.AuthService.RefreshToken(refreshToken)
	if err != nil {
		var userErr *types.MyError
		if errors.As(err, &userErr) {
			switch userErr.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": userErr.Message})
				return
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": userErr.Message})
				return
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": userErr.Message})
				return
			}
		}

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userDto, err := dtos.ConvertUserModelToDTO(*user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the access token and user info
	ctx.JSON(http.StatusOK, gin.H{
		"refreshToken": newRefreshToken,
		"accessToken":  newAccessToken,
		"user":         userDto,
	})
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid refresh token header"})
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	if err := ac.AuthService.Logout(refreshToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ac *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token") // Extract the token from query parameters

	if err := ac.UserService.VerifyEmail(token); err != nil {
		log.Println(err.Error())
		var userErr *types.MyError
		if errors.As(err, &userErr) {
			switch userErr.Code {
			case types.USER_NOT_FOUND_ERROR:
				c.JSON(http.StatusNotFound, gin.H{"error": userErr.Message})
				return
			case types.USER_DEACTIVATED_ERROR:
				c.JSON(http.StatusForbidden, gin.H{"error": userErr.Message})
				return
			case types.INVALID_TOKEN_ERROR:
				c.JSON(http.StatusUnauthorized, gin.H{"error": userErr.Message})
				return
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": userErr.Message})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// StartGoogleAuth initiates the Google OAuth2 flow
func (ac *AuthController) StartGoogleAuth(c *gin.Context) {
	// Add the provider name to the request context
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", "google"))

	// Start the OAuth2 flow
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// GoogleCallback handles the callback from Google after OAuth2 login
func (ac *AuthController) GoogleCallback(c *gin.Context) {
	// Complete the OAuth2 flow and get the user's Google profile
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate with Google"})
		return
	}

	// Use the Google profile to create or log in a user in your system
	// For example, check if the user already exists by email
	existingUser, err := ac.UserService.GetUserByEmail(user.Email)
	if err != nil {
		// If the user doesn't exist, create a new user
		newUser := &dtos.RegisterRequest{
			Email:          user.Email,
			Username:       user.Name,
			Password:       "12345678910",
			Bio:            "Please edit me",
			ProfilePicture: user.AvatarURL,
			DateOfBirth:    "2000-01-01",
			Provider:       user.Provider,
		}
		if err := ac.UserService.RegisterUser(newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		createdUser, err := ac.UserService.GetUserByEmail(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}

		existingUser = createdUser
	}

	// Generate JWT tokens for the user
	accessToken, refreshToken, err := ac.AuthService.GenerateToken(existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	userDto, err := dtos.ConvertUserModelToDTO(*existingUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serialize the userDto to JSON
	userDtoJSON, err := json.Marshal(userDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize user data"})
		return
	}

	// Redirect to the frontend callback route with the OAuth data
	frontendCallbackURL := os.Getenv("FRONTEND_URL") + "/auth/callback"
	redirectURL := fmt.Sprintf(
		"%s?accessToken=%s&refreshToken=%s&user=%s",
		frontendCallbackURL,
		accessToken,
		refreshToken,
		url.QueryEscape(string(userDtoJSON)),
	)
	c.Redirect(http.StatusFound, redirectURL)
}
