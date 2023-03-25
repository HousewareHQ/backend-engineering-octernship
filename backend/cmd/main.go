package main

import (
	"fmt"

	"github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/database"
	authRouter "github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/routes/authRouter"
	orgRouter "github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/routes/orgRouter"
	"github.com/gofiber/fiber"
)

func main() {
	var PORT = database.EnvMap["NGINX_PORT"]
	fmt.Println("Starting server on port:", PORT)

	app := fiber.New()

	// ! routes
	authRouter.SetUp(app)
	orgRouter.SetUp(app)

	app.Listen(PORT)
}
