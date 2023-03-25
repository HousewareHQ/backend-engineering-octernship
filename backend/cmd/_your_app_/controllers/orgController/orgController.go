package orgController

import (
	"net/http"

	"github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/database"
	organisations "github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/models/organisation"
	users "github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/models/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
)

var SecretKey = database.EnvMap["SECRET_KEY"]

func CreateOrg(ctx *fiber.Ctx) {
	var data map[string]string

	// ! parse json body
	if err := ctx.BodyParser(&data); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid organisation data",
		})
		return
	}

	// ! validate data
	if data["name"] == "" || data["head"] == "" {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid organisation data",
		})
		return
	}

	// ! check if org already exists
	var foundOrg organisations.Organisation
	database.DB.Where("name = ?", data["name"]).First(&foundOrg)
	if foundOrg.Id != 0 {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Organisation already exists",
		})
		return
	}

	// ! create org
	org := organisations.Organisation{
		Name: data["name"],
		Head: data["head"],
	}

	// ! save to db
	database.DB.Create(&org)

	ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Organisation created",
		"data":    org,
	})

}

func GetOrgPeople(ctx *fiber.Ctx) {
	cookie := ctx.Cookies("jwt")

	// ! parse jwt token
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token",
		})
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)
	var user users.User
	database.DB.Where("id = ?", claims.Issuer).First(&user)

	orgId := user.OrgId

	// fetch other users in the same org
	var users []users.User
	database.DB.Where("org_id = ?", orgId).Find(&users)

	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Users in organisation fetched",
		"data":    users,
	})
}

func DeleteOrg(ctx *fiber.Ctx) {
	var data map[string]string
	// ! parse json body
	if err := ctx.BodyParser(&data); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid organisation data",
		})
		return
	}

	// ! validate data
	if data["name"] == "" {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid organisation data",
		})
		return
	}

	// ! check if org exists
	var foundOrg organisations.Organisation
	database.DB.Where("name = ?", data["name"]).First(&foundOrg)

	if foundOrg.Id == 0 {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Organisation does not exist",
		})
		return
	}

	// ! delete org
	database.DB.Delete(&foundOrg)

	ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Organisation deleted",
	})

}

func UpdateOrg(ctx *fiber.Ctx) {
	var data map[string]string
	// ! parse json body
	if err := ctx.BodyParser(&data); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid organisation data",
		})
		return
	}

	// ! validate data
	if data["name"] == "" {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid organisation data",
		})
		return
	}

	// ! check if org exists
	var foundOrg organisations.Organisation
	database.DB.Where("name = ?", data["name"]).First(&foundOrg)

	if foundOrg.Id == 0 {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Organisation does not exist",
		})
		return
	}

	// ! update org
	database.DB.Model(&foundOrg).Updates(organisations.Organisation{
		Name: data["name"],
		Head: data["head"],
	})

	ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Organisation updated",
		"data":    foundOrg,
	})
}
