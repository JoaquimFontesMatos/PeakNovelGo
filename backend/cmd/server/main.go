package main

import (
	"backend/config"
	"backend/internal/auth"
	"backend/internal/controllers"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/routes"
	"backend/internal/services"
	"backend/internal/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default port for local testing
	}

	db := config.ConnectDB(false)

	ttsDir := "./tts-files"
	logDir := "logs"

	// Delete all TTS files
	err := os.RemoveAll(ttsDir)
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
	logRepo := repositories.NewLogRepository(db)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, authRepo)
	novelService := services.NewNovelService(novelRepo)
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
	ttsController := controllers.NewTTSController(ttsService)
	logController := controllers.NewLogController(logFilePath, logService)

	// Set up routes
	routes.SetupRoutes(r, authController, userController, novelController, ttsController, logController)

	fmt.Printf("Server running on port %s\n", port)
	err = r.Run(":" + port)
	if err != nil {
		return
	} // Start server on port 8080
}
