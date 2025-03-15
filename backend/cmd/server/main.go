package main

import (
	"backend/config"
	_ "backend/docs"
	"backend/internal/auth"
	"backend/internal/controllers"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/routes"
	"backend/internal/services"
	"backend/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

// gin-swagger middleware
// swagger embed files

//	@title			PeakNovelGo API
//	@version		1.0.12
//	@description	This is the api of the peaknovelgo website.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8081
//	@BasePath	/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

// @securityDefinitions.apikey RefreshBearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the refresh token.

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default port for local testing
	}

	db := config.ConnectDB(false)

	ttsDir := "./tts-files"
	logDir := "logs"

	// Delete all TTS files
	err = os.RemoveAll(ttsDir)
	if err != nil {
		log.Println(err.Error())
	}

	// Ensure the logs directory exists.
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Error creating logs directory: %v", err)
	}

	// Use an absolute path to the log file.
	logFilePath := filepath.Join(logDir, "app.log")

	r := gin.Default()
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	novelRepo := repositories.NewNovelRepository(db)
	chapterRepo := repositories.NewChapterRepository(db)
	bookmarkRepo := repositories.NewBookmarkRepository(db)
	logRepo := repositories.NewLogRepository(db)

	scriptExecutor := &utils.RealScriptExecutor{}

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, authRepo)
	novelService := services.NewNovelService(novelRepo, scriptExecutor)
	chapterService := services.NewChapterService(chapterRepo)
	bookmarkService := services.NewBookmarkService(bookmarkRepo)
	logService := services.NewLogService(logRepo)

	ttsService := &services.TTSService{
		OutputDir: ttsDir,
	}

	_, err = userRepo.GetUserByEmail("admin@example.com")
	if err != nil {
		// Hash the user's password
		hashedPassword, err := utils.HashPassword(os.Getenv("ADMIN_SECRET"))
		if err != nil {
			log.Println(err)
			return
		}

		err = userRepo.CreateUser(&models.User{
			Email:          "admin@example.com",
			Password:       hashedPassword,
			Roles:          "admin",
			Username:       "Admin",
			ProfilePicture: "https://cdn3.iconfinder.com/data/icons/user-group-black/100/user-process-512.png",
			Bio:            "Admin user",
		})
		if err != nil {
			log.Println(err)
			return
		}
	}

	auth.NewAuth()

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService, userService)
	novelController := controllers.NewNovelController(novelService)
	bookmarkController := controllers.NewBookmarkController(bookmarkService)
	chapterController := controllers.NewChapterController(chapterService, novelService)
	ttsController := controllers.NewTTSController(ttsService)
	logController := controllers.NewLogController(logFilePath, logService)

	// Set up middleware
	middleware := middleware.NewMiddleware(userService)

	// Set up routes
	routes.SetupRoutes(r, authController, userController, novelController, bookmarkController, chapterController, ttsController, logController, middleware)

	fmt.Printf("Server running on port %s\n", port)
	err = r.Run(":" + port)
	if err != nil {
		return
	}
}
