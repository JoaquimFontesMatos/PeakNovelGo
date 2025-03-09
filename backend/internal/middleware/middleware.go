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

// Middleware struct represents a middleware component.
//
// Fields:
//   - userService (interfaces.UserServiceInterface): An interface for user service operations.  Used by the middleware for user-related logic.
type Middleware struct {
	userService interfaces.UserServiceInterface
}

// NewMiddleware creates a new middleware instance.
//
// Parameters:
//   - userService (interfaces.UserServiceInterface): The user service to be used by the middleware.
//
// Returns:
//   - *Middleware: A pointer to the newly created middleware instance.
func NewMiddleware(userService interfaces.UserServiceInterface) *Middleware {
	return &Middleware{
		userService: userService,
	}
}

// AuthMiddleware is a Gin middleware function that authenticates users based on a JWT (JSON Web Token) in the Authorization
// header.
// It retrieves the JWT, verifies its signature and expiration, extracts user ID from claims, fetches user details from
// the database, and sets the user information in the Gin context.
// If authentication fails at any stage, it returns an appropriate error response and aborts the request.
//
// Parameters:
//   - m (*Middleware): The middleware instance containing dependencies like userService.
//
// Returns:
//   - gin.HandlerFunc: A Gin middleware handler function.
//
// Error types:
//   - error: Various errors during JWT verification, user fetching from the database, or invalid token claims.  These
//  result in HTTP status codes like 401 (Unauthorized).
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

		// Extract the token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Validate token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
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

		c.Next()
	}
}

// RefreshTokenMiddleware is a middleware function that validates a refresh token from the Authorization header.
//
// It extracts the refresh token from the "Authorization" header, verifies its signature using the SECRET_KEY environment variable,
// and if valid, allows the request to proceed.  If the token is missing, invalid, or verification fails, it returns a 401 Unauthorized error.
//
// Parameters:
//   -  (none):
//
// Returns:
//   - gin.HandlerFunc: A Gin middleware handler function.
//
// Error types:
//   - Error:  Returns a 401 Unauthorized error with an "error" field indicating "Missing or invalid refresh token header" or "Invalid refresh token" if token validation fails.  A fatal error occurs if the SECRET_KEY environment variable is not set.
func (m *Middleware) RefreshTokenMiddleware() gin.HandlerFunc {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY not set in environment variables")
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid refresh token header"})
			c.Abort()
			return
		}

		// Extract the refresh token from the header
		refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the refresh token
		_, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
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

// PermissionMiddleware checks if the user has permission to access a resource and perform a specific action.
//
// It retrieves the user from the Gin context, type asserts it to models.User, and then uses the permissions package
// to check if the user has the required permission.  If the user lacks permission, a 403 Forbidden response is returned.
// If the user is not found or the data is invalid, a 401 Unauthorized response is returned.
//
// Parameters:
//  - resource (string): The resource being accessed (e.g., "users").
//  - action (string): The action being performed (e.g., "read", "write", "delete").
//
// Returns:
//   - gin.HandlerFunc: A Gin middleware function that checks permissions.
//
// Error types:
//   - 401 Unauthorized: Returned if the user is not authenticated or the user data is invalid.
//   - 403 Forbidden: Returned if the user does not have permission to access the resource and perform the action.
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

