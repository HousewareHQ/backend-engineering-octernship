package routes

import (
	"github.com/gin-gonic/gin"

	controllers "github.com/HousewareHQ/backend-engineering-octernship/api/server/controllers"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/middlewares"
)

func UserRoutes(routeUrl *gin.Engine) {
	// Middlware to verify authentication before accessing following routes
	routeUrl.Use(middlewares.Authenticate()) //User authorization level

	routeUrl.GET("/users", controllers.GetAllUsers())

}
