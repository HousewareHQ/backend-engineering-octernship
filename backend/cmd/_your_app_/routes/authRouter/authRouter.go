package authRouter

import (
	authController "github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/controllers/authController"
	"github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/middleware/isAuth"
	"github.com/gofiber/fiber"
)

func SetUp(app *fiber.App) {
	// ! admin routes
	app.Post("/api/auth", isAuth.IsAdminCheck, authController.SignUp)
	app.Delete("/api/auth", isAuth.IsAdminCheck, authController.DeleteUser)
	// ! public routes
	app.Post("/api/auth/login", authController.Login)
	app.Get("/api/auth/user", authController.GetUser)
	app.Get("api/auth/logout", authController.Logout)
	app.Get("api/auth/refresh", authController.RefreshToken)
}
