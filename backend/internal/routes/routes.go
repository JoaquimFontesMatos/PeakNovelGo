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
		user.GET("/:id", userController.HandleGetUser)
		user.GET("/email/:email", middleware.AuthMiddleware(), userController.HandleGetUserByEmail)
		user.GET("/username/:username", middleware.AuthMiddleware(), userController.HandleGetUserByUsername)
		user.PUT("/:id/password", middleware.AuthMiddleware(), userController.UpdatePassword)
		user.PUT("/:id/email", middleware.AuthMiddleware(), userController.UpdateEmail)
		user.PUT("/:id/fields", middleware.AuthMiddleware(), userController.UpdateUserFields)
		user.DELETE("/:id", middleware.AuthMiddleware(), userController.HandleDeleteUser)
	}

	novel := r.Group("/novels")
	{
		novel.POST("/", novelController.HandleImportNovel)
		novel.GET("/authors/:author_id", novelController.GetNovelsByAuthorID)
		novel.GET("/genres/:genre_id", novelController.GetNovelsByGenreID)
		novel.GET("/tags/:tag_id", novelController.GetNovelsByTagID)
		novel.GET("/:novel_id", novelController.GetNovelByID)

		novel := r.Group("/novels/chapters")
		{
			novel.POST("/:novel_id", novelController.HandleImportChaptersZip)
			novel.GET("/:novel_id", novelController.GetChaptersByNovelID)
			novel.GET("/chapter/:chapter_id", novelController.GetChapterByID)
		}
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
