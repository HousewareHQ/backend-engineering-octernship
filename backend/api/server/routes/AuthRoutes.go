package routes

import (
	controllers "github.com/HousewareHQ/backend-engineering-octernship/api/server/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(routeUrl *gin.Engine) {
	routeUrl.POST("/users/login", controllers.Login())

}
