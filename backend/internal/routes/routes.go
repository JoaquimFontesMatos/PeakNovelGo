package routes

import (
	"backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *controllers.UserController) *gin.Engine {
	router := gin.Default()

	// Define user routes
	router.POST("/users", userController.HandleCreateUser)
	router.GET("/users/:id", userController.HandleGetUser)

	return router
}
