package main

import (
	"backend/config"
	"backend/internal/controllers"
	"backend/internal/repositories"
	"backend/internal/routes"
	"backend/internal/services"
)

func main() {
	db := config.ConnectDB(false)
	userRepo := repositories.UserRepository{DB: db}
	userService := services.UserService{Repo: &userRepo}
	userController := controllers.UserController{Service: &userService}

	router := routes.SetupRouter(&userController)
	router.Run(":8080") // Start server on port 8080
}
