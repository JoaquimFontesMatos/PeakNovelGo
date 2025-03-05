package middleware

import (
	"backend/internal/models"
	"backend/internal/permissions"
	"backend/internal/services/interfaces"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Middleware struct holds dependencies for the middleware
type Middleware struct {
	userService interfaces.UserServiceInterface
}

// NewMiddleware creates a new instance of the Middleware
func NewMiddleware(userService interfaces.UserServiceInterface) *Middleware {
	return &Middleware{
		userService: userService,
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY not set in environment variables")
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		// Extract and parse the token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		// Handle token parsing errors
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		// Validate token claims (e.g., expiration)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expiration, ok := claims["exp"].(float64)
			if !ok || time.Unix(int64(expiration), 0).Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}

			// Extract user information from claims
			userID, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}

			// Fetch the user from the database (or wherever you store user data)
			user, err := m.userService.GetUser(uint(userID))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				c.Abort()
				return
			}

			// Set the user in the Gin context
			c.Set("user", user)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}

// RefreshTokenMiddleware returns a Gin middleware for refreshing tokens
func (m *Middleware) RefreshTokenMiddleware() gin.HandlerFunc {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY not set in environment variables")
	}

	return func(c *gin.Context) {
		// Check for refresh token in cookies
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not provided"})
			c.Abort()
			return
		}

		// Parse and validate the refresh token
		_, err = jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PermissionMiddleware returns a Gin middleware to check permissions for a specific resource and action
func (m *Middleware) PermissionMiddleware(resource string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user from the Gin context (assuming you set it in a previous middleware, e.g., AuthMiddleware)
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Type assert the user to models.User
		userModel, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
			c.Abort()
			return
		}

		// Check if the user has permission
		if !permissions.HasPermission(*userModel, resource, action, nil) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// Proceed to the next handler if the user has permission
		c.Next()
	}
}
