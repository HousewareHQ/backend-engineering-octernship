package main

import (
	"os"

	"github.com/HousewareHQ/backend-engineering-octernship/api/server/middlewares"
	routes "github.com/HousewareHQ/backend-engineering-octernship/api/server/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port, portExists := os.LookupEnv("PORT")

	if !portExists {
		port = "8080"
	}
	router := gin.New()
	router.Use(middlewares.CORSMiddleware())
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.AdminRoutes(router)

	router.Run(":" + port)
}
