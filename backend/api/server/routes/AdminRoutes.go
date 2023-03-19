package routes

import (
	controllers "github.com/HousewareHQ/backend-engineering-octernship/api/server/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(routeUrl *gin.Engine) {
	routeUrl.POST("/users/create-user", controllers.CreateUser())
	routeUrl.DELETE("/users/delete/:uid", controllers.DeleteUser())

}
