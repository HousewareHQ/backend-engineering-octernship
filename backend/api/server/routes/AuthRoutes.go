package routes

import (
	controllers "github.com/HousewareHQ/backend-engineering-octernship/api/server/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(routeUrl *gin.Engine) {
	routeUrl.POST("/login", controllers.Login())   //login user
	routeUrl.POST("/logout", controllers.Logout()) //logout user

}
