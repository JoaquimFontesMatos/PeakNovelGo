package main

import (
	"backend/config"
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

// main is the entry point of the application. It initializes the database connection, sets up routes, and starts the server.
// It also handles creation of the admin user if it doesn't exist.
//
// Parameters:
//   - None
//
// Returns:
//   - None
//
// Error types:
//   - error:  Various errors can occur during database connection, environment variable loading, file system operations,
//     or server startup. These are logged but not explicitly returned.
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
