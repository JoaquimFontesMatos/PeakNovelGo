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
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	authService := services.NewAuthService(*userRepo)
	authController := controllers.NewAuthController(*authService, *userService)

	// Set up routes
	routes.SetupRoutes(r, authController, userController)

	fmt.Printf("Server running on port %s\n", port)
	r.Run(":8080") // Start server on port 8080
}
