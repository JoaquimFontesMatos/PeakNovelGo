package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"backend/internal/dtos"
	"backend/internal/services/interfaces"
	"backend/internal/utils"
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
		utils.HandleError(c, err)
		return
	}

	// Call the user service to create a new user
	if err := ac.UserService.RegisterUser(&registerRequest); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user login and returns a JWT token.
//
// @Summary User login
// @Description Authenticates a user and returns a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body dtos.LoginRequest true "Login credentials"
// @Success 200 {object} dtos.AuthSession
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /auth/login [post]
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
		utils.HandleError(c, err)
		return
	}

	// Generate JWT tokens (access token and refresh token)
	accessToken, refreshToken, err := ac.AuthService.GenerateToken(user)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	userDto, err := dtos.ConvertUserModelToDTO(*user)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	authSession := dtos.AuthSession{
		User:         userDto,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	// Respond with the access token and user info
	c.JSON(http.StatusOK, authSession)
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
		utils.HandleError(ctx, err)
		return
	}

	userDto, err := dtos.ConvertUserModelToDTO(*user)
	if err != nil {
		utils.HandleError(ctx, err)
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
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ac *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token") // Extract the token from query parameters

	if err := ac.UserService.VerifyEmail(token); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// StartOAuth2 initiates the OAuth2 flow
func (ac *AuthController) StartOAuth2(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	session, _ := gothic.Store.Get(c.Request, gothic.SessionName)
	session.Values["provider"] = provider
	session.Save(c.Request, c.Writer)

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// Oauth2Callback handles the callback from the provider after OAuth2 login
func (ac *AuthController) Oauth2Callback(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	// Use the profile to create or log in a user in your system
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
