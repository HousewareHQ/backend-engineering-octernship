package authController

import (
	"net/http"
	"strconv"
	"time"

	"github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/database"
	tokens "github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/models/jwt_token"
	organisations "github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/models/organisation"
	users "github.com/HousewareHQ/backend-engineering-octernship/cmd/_your_app_/models/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
)

// todo - move to .env
const EXPIRATION_TIME = 10 * time.Second
const REFRESH_TOKEN_TIME = 1 * time.Minute

var SecretKey = database.EnvMap["SECRET_KEY"]

func SignUp(ctx *fiber.Ctx) {
	var data map[string]string

	// ! parse json body
	if err := ctx.BodyParser(&data); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user data",
		})
		return
	}

	// ! validate data
	if data["username"] == "" || data["password"] == "" {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user data",
		})
		return
	}

	// ! check if user already exists
	var foundUser users.User
	database.DB.Where("username = ?", data["username"]).First(&foundUser)
	if foundUser.Id != 0 {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "User already exists",
		})
		return
	}

	// ! hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 8)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to hash password",
		})
		return
	}

	// typecast string to uint
	orgId, _ := strconv.ParseUint(data["org_id"], 0, 0)

	// ! get Organisation object
	var org organisations.Organisation
	database.DB.Where("id = ?", orgId).First(&org)

	// check if such org exists
	if org.Id == 0 {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "No such organisation",
		})
		return
	}

	// ! create user
	user := users.User{
		Username:     data["username"],
		Password:     hashedPass,
		OrgId:        uint(orgId),
		Organisation: org,
	}

	// ! save user to db
	database.DB.Create(&user)

	// ! send user object to client in form of json
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User created",
		"data":    user,
	})
}

func Login(ctx *fiber.Ctx) {
	var data map[string]string

	// ! parse json body
	if err := ctx.BodyParser(&data); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user data",
		})
	}

	// ! validate data
	if data["username"] == "" || data["password"] == "" {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user data",
		})
		return
	}

	// ! check if user exists
	var user users.User
	database.DB.Where("username = ?", data["username"]).First(&user)
	if user.Id == 0 {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "No such user",
		})
		return
	}

	// ! check if password is correct
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Incorrect password",
		})
		return
	}

	// ! create jwt token that also has username and user id that has an expiry time of 24 hours
	claims := &tokens.Claims{
		Id:       user.Id,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		OrgId:    user.OrgId,
		StandardClaims: jwt.StandardClaims{
			// ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			ExpiresAt: time.Now().Add(EXPIRATION_TIME).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to create token",
		})
		return
	}

	// ! store token in cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(EXPIRATION_TIME),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logged in",
	})
}

func GetUser(ctx *fiber.Ctx) {
	cookie := ctx.Cookies("jwt")

	// ! parse jwt token
	token, err := jwt.ParseWithClaims(cookie, &tokens.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid token",
		})
		return
	}

	claims := token.Claims.(*tokens.Claims)
	var user users.User

	// ! user obj is stored in claims
	user.Id = claims.Id
	user.Username = claims.Username
	user.OrgId = claims.OrgId
	user.IsAdmin = claims.IsAdmin

	if user.Id == 0 {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "No such user",
		})
		return
	}

	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User found",
		"data":    user,
	})
}

func RefreshToken(ctx *fiber.Ctx) {
	cookie := ctx.Cookies("jwt")

	// ! parse jwt token
	token, err := jwt.ParseWithClaims(cookie, &tokens.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token",
		})
		return
	}

	claims := token.Claims.(*tokens.Claims)
	// ! update expiry time
	claims.StandardClaims.ExpiresAt = time.Now().Add(REFRESH_TOKEN_TIME).Unix()

	// ! create new token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString([]byte(SecretKey))

	if err != nil {
		ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to create token",
		})
		return
	}

	// ! store token in cookie
	newCookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(REFRESH_TOKEN_TIME),
		HTTPOnly: true,
	}

	ctx.Cookie(&newCookie)
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Token refreshed",
	})
}

func Logout(ctx *fiber.Ctx) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User logged out successfully",
	})
}

func DeleteUser(ctx *fiber.Ctx) {
	// ! parse json body
	var data map[string]string
	if err := ctx.BodyParser(&data); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user data",
		})
	}

	// ! validate data
	if data["username"] == "" {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user data",
		})
		return
	}

	// ! delete user
	database.DB.Where("username = ?", data["username"]).Delete(&users.User{})
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted",
	})
}
