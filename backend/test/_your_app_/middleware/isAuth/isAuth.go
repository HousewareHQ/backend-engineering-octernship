package isAuth

import (
	"net/http"

	"github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/database"
	tokens "github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/models/jwt_token"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
)

var SecretKey = database.EnvMap["SECRET_KEY"]

func IsAdminCheck(ctx *fiber.Ctx) {
	cookie := ctx.Cookies("jwt")

	// ! parse jwt token
	token, err := jwt.ParseWithClaims(cookie, &tokens.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token",
		})
		return
	}

	claims := token.Claims.(*tokens.Claims)

	if claims.Id == 0 || !claims.IsAdmin {
		ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorised Access",
		})
		return
	}

	ctx.Next()
}
