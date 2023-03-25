package orgRouter

import (
	orgController "github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/controllers/orgController"
	"github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/middleware/isAuth"
	"github.com/gofiber/fiber"
)

func SetUp(app *fiber.App) {
	// ! admin routes
	app.Post("/api/org", isAuth.IsAdminCheck, orgController.CreateOrg)
	app.Delete("/api/org", isAuth.IsAdminCheck, orgController.DeleteOrg)
	app.Patch("api/org", isAuth.IsAdminCheck, orgController.UpdateOrg)
	// ! public routes
	app.Get("/api/org", orgController.GetOrgPeople)
}
