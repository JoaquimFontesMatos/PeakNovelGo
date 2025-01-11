package controllers

import (
	"errors"
	"net/http"

	"backend/internal/dtos"
	"backend/internal/services/interfaces"
	"backend/internal/types"
	"backend/internal/utils"
	"backend/internal/validators"

	"github.com/gin-gonic/gin"
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
			case types.CONFLICT_ERROR:
				c.JSON(http.StatusConflict, gin.H{"error": error.Message})
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.INTERNAL_SERVER_ERROR:
				c.JSON(http.StatusInternalServerError, gin.H{"error": error.Message})
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
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

	// Respond with both tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	refreshToken, err := utils.ExtractToken(authHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
		return
	}

	// Call the service to refresh the token
	newAccessToken, newRefreshToken, err := c.AuthService.RefreshToken(refreshToken)
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

	// Send the new tokens back to the client
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	refreshToken, err := utils.ExtractToken(authHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
		return
	}

	if err := ac.AuthService.Logout(refreshToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ac *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token") // Extract the token from query parameters

	if err := ac.UserService.VerifyEmail(token); err != nil {
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
