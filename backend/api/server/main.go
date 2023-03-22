package main

import (
	"os"

	"github.com/HousewareHQ/backend-engineering-octernship/api/server/docs"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/middlewares"
	routes "github.com/HousewareHQ/backend-engineering-octernship/api/server/routes"
	"github.com/gin-gonic/gin"
	swaggeFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @contact.name   	Pranay Payal
// @contact.url    	https://www.linkedin.com/in/pranay-payal-b6b0161b1/
// @contact.email  	kewinlee123@gmail.com
// @license.name  	Apache 2.0
// @license.url   	http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiToken
// @in header
// @name Authorization

func main() {
	//Swagger Meta information
	docs.SwaggerInfo.Title = "Organization's API"
	docs.SwaggerInfo.Description = "Authentication-Authorization Service API.<br>React client at http://localhost:3000"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	port, portExists := os.LookupEnv("PORT") //GET PORT from .env file

	if !portExists {
		port = "8080" //if fails assign 8080
	}
	router := gin.New()                                  //Creating Gin Router
	router.Use(gin.Logger())                             //gin logger-middlware
	router.Use(middlewares.PreflightRequestMiddleware()) //Preflight Request middleware

	baseRouter := router.Group("/api/v1") //Adding base path

	//Adding Swagger
	baseRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggeFiles.Handler))

	routes.AuthRoutes(baseRouter)  //AUTHENTICATION routes
	routes.UserRoutes(baseRouter)  //USER routes
	routes.AdminRoutes(baseRouter) //ADMIN routes

	router.Run(":" + port) //START SERVER
}
