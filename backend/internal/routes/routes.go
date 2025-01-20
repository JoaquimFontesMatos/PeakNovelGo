package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	novelController *controllers.NovelController,
	ttsController *controllers.TTSController) {

	// Serve TTS files
	r.Static("/tts-files", "./tts-files")

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
		novel.POST("/:novel_updates_id", novelController.HandleImportNovel)
		novel.GET("/", novelController.GetNovels)
		novel.GET("/authors/:author_name", novelController.GetNovelsByAuthorName)
		novel.GET("/genres/:genre_name", novelController.GetNovelsByGenreName)
		novel.GET("/tags/:tag_name", novelController.GetNovelsByTagName)
		novel.GET("/:novel_id", novelController.GetNovelByID)
		novel.GET("/title/:title", novelController.GetNovelByUpdatesID)

		chapters := r.Group("/novels/chapters")
		{
			chapters.POST("/:novel_id", novelController.HandleImportChaptersZip)
			chapters.GET("/:novel_id", novelController.GetChaptersByNovelID)
			chapters.GET("/chapter/:chapter_id", novelController.GetChapterByID)
			chapters.GET("/novel/:novel_title/chapter/:chapter_no", novelController.GetChapterByNovelUpdatesIDAndChapterNo)
			chapters.GET("/novel/:novel_title/chapters", novelController.GetChaptersByNovelUpdatesID)
		}

		bookmarked := r.Group("/novels/bookmarked")
		{
			bookmarked.POST("/", novelController.CreateBookmarkedNovel)
			bookmarked.PUT("/", novelController.UpdateBookmarkedNovel)
			bookmarked.GET("/:user_id", novelController.GetBookmarkedNovelsByUserID)
			bookmarked.GET("/user/:user_id/novel/:novel_id", novelController.GetBookmarkedNovelByUserIDAndNovelID)
		}

		tts := r.Group("/novels/tts")
		{
			tts.POST("/", ttsController.GenerateTTS)
			tts.GET("/voices", ttsController.GetVoices)
		}
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
