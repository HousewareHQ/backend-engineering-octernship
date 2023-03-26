package routes

import (
	controller "github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/controllers"
	middleware "github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.Use(middleware.TokenAuth())
	incomingRoutes.GET("users", controller.GetUsers())	// this orgId should match that of the user/admin, then only can send this request
	incomingRoutes.POST("logout", controller.Logout())
}


