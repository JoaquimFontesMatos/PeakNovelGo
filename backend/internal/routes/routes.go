package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all the routes for the application.
//
// Parameters:
//   - r (*gin.Engine): The Gin engine to register routes on.
//   - authController (*controllers.AuthController): The authentication controller.
//   - userController (*controllers.UserController): The user controller.
//   - novelController (*controllers.NovelController): The novel controller.
//   - bookmarkController (*controllers.BookmarkController): The bookmark controller.
//   - chapterController (*controllers.ChapterController): The chapter controller.
//   - ttsController (*controllers.TTSController): The TTS controller.
//   - logController (*controllers.LogController): The log controller.
//   - middleware (*middleware.Middleware): The middleware to use for authentication and authorization.
func SetupRoutes(r *gin.Engine,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	novelController *controllers.NovelController,
	bookmarkController *controllers.BookmarkController,
	chapterController *controllers.ChapterController,
	ttsController *controllers.TTSController,
	logController *controllers.LogController,
	middleware *middleware.Middleware) {

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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
		novel.POST("/:novel_updates_id", middleware.AuthMiddleware(), middleware.PermissionMiddleware("novels", "create"), novelController.HandleImportNovelByNovelUpdatesID)
		novel.GET("/", novelController.GetNovels)
		novel.GET("/authors/:author_name", novelController.GetNovelsByAuthorName)
		novel.GET("/genres/:genre_name", novelController.GetNovelsByGenreName)
		novel.GET("/tags/:tag_name", novelController.GetNovelsByTagName)
		novel.GET("/:novel_id", novelController.GetNovelByID)
		novel.GET("/title/:title", novelController.GetNovelByUpdatesID)

		chapters := novel.Group("/chapters")
		{
			chapters.GET("/:novel_id/scrape", chapterController.HandleImportChapters)
			chapters.GET("/novel/:novel_title/chapter/:chapter_no", chapterController.GetChapterByNovelUpdatesIDAndChapterNo)
			chapters.GET("/novel/:novel_title/chapters", chapterController.GetChaptersByNovelUpdatesID)
		}

		bookmarked := novel.Group("/bookmarked")
		{
			bookmarked.POST("/", middleware.AuthMiddleware(), bookmarkController.CreateBookmark)
			bookmarked.PUT("/", middleware.AuthMiddleware(), bookmarkController.UpdateBookmark)
			bookmarked.GET("/:user_id", middleware.AuthMiddleware(), bookmarkController.GetBookmarkedNovelsByUserID)
			bookmarked.GET("/user/:user_id/novel/:novel_id", middleware.AuthMiddleware(), bookmarkController.GetBookmarkByUserIDAndNovelID)
			bookmarked.DELETE("/user/:user_id/novel/:novel_id", middleware.AuthMiddleware(), bookmarkController.UnbookmarkNovel)
		}

		tts := novel.Group("/tts")
		{
			tts.POST("/", middleware.AuthMiddleware(), ttsController.GenerateTTS)
			tts.GET("/voices", ttsController.GetVoices)
		}
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
