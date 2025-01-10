package main

import (
	"backend/config"
	"backend/internal/controllers"
	"backend/internal/repositories"
	"backend/internal/routes"
	"backend/internal/services"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for local testing
	}

	db := config.ConnectDB(false)

	r := gin.Default()
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	novelRepo := repositories.NewNovelRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, authRepo)
	novelService := services.NewNovelService(novelRepo)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService, userService)
	novelController := controllers.NewNovelController(novelService)

	// Set up routes
	routes.SetupRoutes(r, authController, userController, novelController)

	fmt.Printf("Server running on port %s\n", port)
	r.Run(":8080") // Start server on port 8080
}
