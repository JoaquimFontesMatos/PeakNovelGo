package routes

import (
	"backend/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authController *controllers.AuthController, userController *controllers.UserController) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh-token", authController.RefreshToken)
		auth.GET("/verify-email", authController.VerifyEmail)
	}
	user := r.Group("/user")
	{
		user.GET("/users/:id", userController.HandleGetUser)
		user.GET("/users/email/:email", userController.HandleGetUserByEmail)
		user.GET("/users/username/:username", userController.HandleGetUserByUsername)
		user.PUT("/users/:id", userController.HandleUpdateUser)
		user.DELETE("/users/:id", userController.HandleDeleteUser)
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
