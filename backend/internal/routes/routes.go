package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	novelController *controllers.NovelController,
	ttsController *controllers.TTSController,
	logController *controllers.LogController) {

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusOK)
			return
		}
		c.Next()
	})

	// Serve TTS files
	r.Static("/tts-files", "./tts-files")

	r.StaticFile("/", "./static/index.html")

	log := r.Group("/log")
	{
		log.POST("/", logController.SaveLog)
		log.GET("/", logController.GetLogs)
		log.GET("/:level", logController.GetLogsByLevel)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh-token", middleware.RefreshTokenMiddleware(), authController.RefreshToken)
		auth.GET("/verify-email", authController.VerifyEmail)
		auth.POST("/logout", middleware.RefreshTokenMiddleware(), authController.Logout)
		auth.GET("/google", authController.StartGoogleAuth)
		auth.GET("/google/callback", authController.GoogleCallback)
	}

	user := r.Group("/user")
	{
		user.GET("/:id", userController.HandleGetUser)
		user.GET("/email/:email", middleware.AuthMiddleware(), userController.HandleGetUserByEmail)
		user.GET("/username/:username", userController.HandleGetUserByUsername)
		user.PUT("/:id/password", middleware.AuthMiddleware(), userController.UpdatePassword)
		user.PUT("/:id/email", middleware.AuthMiddleware(), userController.UpdateEmail)
		user.PUT("/:id/fields", middleware.AuthMiddleware(), userController.UpdateUserFields)
		user.DELETE("/:id", middleware.AuthMiddleware(), userController.HandleDeleteUser)
	}

	novel := r.Group("/novels")
	{
		novel.POST("/", middleware.AuthMiddleware(), novelController.HandleImportNovel)
		novel.POST("/:novel_updates_id", middleware.AuthMiddleware(), novelController.HandleImportNovelByNovelUpdatesID)
		novel.GET("/", novelController.GetNovels)
		novel.GET("/authors/:author_name", novelController.GetNovelsByAuthorName)
		novel.GET("/genres/:genre_name", novelController.GetNovelsByGenreName)
		novel.GET("/tags/:tag_name", novelController.GetNovelsByTagName)
		novel.GET("/:novel_id", novelController.GetNovelByID)
		novel.GET("/title/:title", novelController.GetNovelByUpdatesID)

		chapters := novel.Group("/chapters")
		{
			chapters.POST("/:novel_id/scrape", novelController.HandleImportChapters)
			chapters.POST("/:novel_id", middleware.AuthMiddleware(), novelController.HandleImportChaptersZip)
			chapters.GET("/:novel_id", novelController.GetChaptersByNovelID)
			chapters.GET("/chapter/:chapter_id", novelController.GetChapterByID)
			chapters.GET("/novel/:novel_title/chapter/:chapter_no", novelController.GetChapterByNovelUpdatesIDAndChapterNo)
			chapters.GET("/novel/:novel_title/chapters", novelController.GetChaptersByNovelUpdatesID)
		}

		bookmarked := novel.Group("/bookmarked")
		{
			bookmarked.POST("/", middleware.AuthMiddleware(), novelController.CreateBookmarkedNovel)
			bookmarked.PUT("/", middleware.AuthMiddleware(), novelController.UpdateBookmarkedNovel)
			bookmarked.GET("/:user_id", middleware.AuthMiddleware(), novelController.GetBookmarkedNovelsByUserID)
			bookmarked.GET("/user/:user_id/novel/:novel_id", middleware.AuthMiddleware(), novelController.GetBookmarkedNovelByUserIDAndNovelID)
			bookmarked.DELETE("/user/:user_id/novel/:novel_id", middleware.AuthMiddleware(), novelController.UnbookmarkNovel)
		}

		tts := novel.Group("/tts")
		{
			tts.POST("/", middleware.AuthMiddleware(), ttsController.GenerateTTS)
			tts.GET("/voices", middleware.AuthMiddleware(), ttsController.GetVoices)
		}
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
