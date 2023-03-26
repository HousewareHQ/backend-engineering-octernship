package routes

import (
	controllers "github.com/HousewareHQ/backend-engineering-octernship/api/server/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(routeUrl *gin.RouterGroup) {
	routeUrl.POST("/users/create-user", controllers.CreateUser())   //Create new user in same org. as current user
	routeUrl.DELETE("/users/delete/:uid", controllers.DeleteUser()) //Delete user by id,if exists in current user's organization

}
