package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authController *controllers.AuthController, userController *controllers.UserController, novelController *controllers.NovelController) {
	r.StaticFile("/", "./static/index.html")

	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh-token", middleware.AuthMiddleware(), authController.RefreshToken)
		auth.GET("/verify-email", authController.VerifyEmail)
	}
	user := r.Group("/user")
	{
		user.GET("/users/:id", userController.HandleGetUser)
		user.GET("/users/email/:email", middleware.AuthMiddleware(), userController.HandleGetUserByEmail)
		user.GET("/users/username/:username", middleware.AuthMiddleware(), userController.HandleGetUserByUsername)
		user.PUT("/users/:id/password", middleware.AuthMiddleware(), userController.UpdatePassword)
		user.PUT("/users/:id/email", middleware.AuthMiddleware(), userController.UpdateEmail)
		user.PUT("/users/:id/fields", middleware.AuthMiddleware(), userController.UpdateUserFields)
		user.DELETE("/users/:id", middleware.AuthMiddleware(), userController.HandleDeleteUser)
	}

	novel := r.Group("/novels")
	{
		novel.POST("/novel", novelController.HandleImportNovel)
		novel.POST("/chapters/:novel_id", novelController.HandleImportChaptersZip)
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
