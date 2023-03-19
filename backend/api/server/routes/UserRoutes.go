package routes

import (
	"github.com/gin-gonic/gin"

	controllers "github.com/HousewareHQ/backend-engineering-octernship/api/server/controllers"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/middlewares"
)

func UserRoutes(routeUrl *gin.Engine) {
	// Middlware to verify authentication before accessing following routes
	routeUrl.Use(middlewares.SessionAuthentication()) //whether a session valid or expired
	routeUrl.Use(middlewares.Authorization())         //User authorization level using token

	routeUrl.GET("/users", controllers.GetAllUsers()) //get all users

}
