package main

import (
    "backend/config"
    "backend/internal/controllers"
    "backend/internal/repositories"
    "backend/internal/routes"
    "backend/internal/services"
    "fmt"
    "log"
    "os"

    "github.com/gin-gonic/gin"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081" // Default port for local testing
    }

    db := config.ConnectDB(false)

    // Delete all TTS files
    err := os.RemoveAll("./tts-files")
    if err != nil {
        log.Printf(err.Error())
    }

    r := gin.Default()
    userRepo := repositories.NewUserRepository(db)
    authRepo := repositories.NewAuthRepository(db)
    novelRepo := repositories.NewNovelRepository(db)
    userService := services.NewUserService(userRepo)
    authService := services.NewAuthService(userRepo, authRepo)
    novelService := services.NewNovelService(novelRepo)

    ttsService := &services.TTSService{
        OutputDir: "./tts-files",
    }

    userController := controllers.NewUserController(userService)
    authController := controllers.NewAuthController(authService, userService)
    novelController := controllers.NewNovelController(novelService)
    ttsController := controllers.NewTTSController(ttsService)

    // Set up routes
    routes.SetupRoutes(r, authController, userController, novelController, ttsController)

    fmt.Printf("Server running on port %s\n", port)
    err = r.Run(":" + port)
    if err != nil {
        return
    } // Start server on port 8080
}
