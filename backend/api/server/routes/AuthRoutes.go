package routes

import (
	controllers "github.com/HousewareHQ/backend-engineering-octernship/api/server/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(routeUrl *gin.RouterGroup) {
	routeUrl.POST("/login", controllers.Login())                       //login user
	routeUrl.POST("/logout", controllers.Logout())                     //logout user
	routeUrl.POST("/refresh-rotate", controllers.RefreshTokenRotate()) //refresh token rotate

	routeUrl.POST("/super/create-user", controllers.CreateSuperUser()) //create first admin to start testing API
}
