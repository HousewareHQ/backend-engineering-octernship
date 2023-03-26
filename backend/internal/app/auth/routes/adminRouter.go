package routes

import (
	controller "github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/controllers"
	"github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/middleware"

	"github.com/gin-gonic/gin"
)


func AdminRoutes(incomingRoutes *gin.Engine) {
	// only admin usertype can access these request - using the authenticate middleware to check user type
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/users", controller.AddUserToOrg())// this orgId should match the orgId of the ADMIN - orgId is in the token
	incomingRoutes.DELETE("/users/:userId", controller.DeleteUserFromOrg())	// this orgId should match the orgId of the ADMIN and only admin can send this request  - middleware checks this
}
	
	




