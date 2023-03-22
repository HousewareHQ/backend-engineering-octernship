package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	AppConstant "github.com/HousewareHQ/backend-engineering-octernship/api/server/constants"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/helpers"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/models"
	DBconnect "github.com/HousewareHQ/backend-engineering-octernship/api/server/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

// Validator and userCollection variable
var userCollection *mongo.Collection = DBconnect.OpenCollection(DBconnect.Client, AppConstant.USER_COLLECTION)
var validate = validator.New()

// Login User 		godoc
// @Summary			Login User
// @Description		Login user and create new session(assigns access-token & refresh-token to user)
// @ID				Authentication
// @Produce 		application/json
// @Consume 		application/json
// @Param			login-request body models.LoginRequestBody true "User Credentials"
// @Tags			Authentication Endpoints
// @Success			200 {object} models.User
// @Failure			500
// @Header 200 {string} Set-Cookie "accesstoken+refreshtoken"
// @Router 			/login [post]
func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var storedUser models.User
		defer cancel()

		//unmarshal encoded-json into struct
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Println(err.Error())
			return
		}

		//Get user document using username as query parameter
		err := userCollection.FindOne(c, bson.M{"username": user.Username}).Decode(&storedUser)

		/* Validations*/
		if user.Username == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "User not found"})
			return
		}
		defer cancel()
		//if document not found
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect username or does not exists"})

		}

		//Validate Password
		if !helpers.VerifyPassword(user.Password, storedUser.Password) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Password does not match"})
			return

		}
		//password matches,User logged in
		/*updating jwt tokens of user*/
		token, refreshedToken := helpers.GenerateTokens(storedUser)
		storedUser.JWTToken = token
		storedUser.RefreshToken = refreshedToken
		//storing tokens on user document in db
		helpers.UpdateTokenOnLogin(token, refreshedToken, storedUser.ID)
		//storing tokens locally in cookies
		ctx.SetCookie("accesstoken", token, int(AppConstant.TOKEN_COOKIE_EXPIRY), "/", "localhost", false, true)
		ctx.SetCookie("refreshtoken", refreshedToken, int(AppConstant.REFRESH_TOKEN_COOKIE_EXPIRY), "/", "localhost", false, true)

		ctx.JSON(http.StatusOK, storedUser)

	}

}

// Logout User 		godoc
// @Summary			Logout User
// @Description		Logout user and Destroys user's session by expiring tokens cookie
// @ID				Logout
// @Produce 		application/json
// @Tags			Authentication Endpoints
// @Success			200
// @Failure			500
// @Header 200 {string} Set-Cookie "(accesstoken+refreshtoken)-expired"
// @Router 			/logout [post]
func Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Throw away current logged in session
		//By removing storedAccess and refresh tokens in cookies
		ctx.SetCookie("accesstoken", "", -1, "/", "localhost", false, true)
		ctx.SetCookie("refreshtoken", "", -1, "/", "localhost", false, true)
		ctx.JSON(http.StatusOK, gin.H{"Ok": "Logout"})
	}
}

// RefreshTokenRotate  godoc
// @Summary			Refresh Token Rotate
// @Description		Returns and assigns new access-token and refresh-token
// @ID				RefreshTokenRotate
// @Produce 		application/json
// @Param			Cookie header string false "Send token using Cookie header<br>(Example:Cookie:refreshtoken=eyJhbGciOiJ..)"
// @Tags			Authentication Endpoints
// @Security 		ApiToken
// @Success			200 {object} models.RefreshRotateResponse
// @Failure			500
// @Header 200 {string} Set-Cookie "accesstoken+refreshtoken"
// @Router 			/refresh-rotate [post]
func RefreshTokenRotate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		//GET refresh token from header
		oldRefreshToken := ctx.GetHeader("Authorization")
		if oldRefreshToken == "" { // if no refresh token provided in header
			var refreshTokenErr error
			//Check if refresh token is present in cookies
			oldRefreshToken, refreshTokenErr = ctx.Cookie("refreshtoken")
			defer cancel()
			if refreshTokenErr != nil { // if fails then show InternalServerError
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": refreshTokenErr.Error()})
				return
			}
		}
		defer cancel()

		//generate new tokens using fetched old refreshtoken
		accessToken, refreshToken, generateTokenErr := helpers.GenerateTokenByRefreshToken(c, oldRefreshToken)
		if generateTokenErr != nil {
			//IF:fail to generate valid token then user is no longer authorized.
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": generateTokenErr.Error()})
			ctx.Abort()
			return

		}
		//Store tokens in cookies
		ctx.SetCookie("accesstoken", accessToken, int(AppConstant.TOKEN_COOKIE_EXPIRY), "/", "localhost", false, true)
		ctx.SetCookie("refreshtoken", refreshToken, int(AppConstant.REFRESH_TOKEN_COOKIE_EXPIRY), "/", "localhost", false, true)
		//send newly generated tokens.
		ctx.JSON(http.StatusOK, gin.H{"accesstoken": accessToken, "refreshtoken": refreshToken})
	}
}

// GetAllUsers 		godoc
// @Summary			Get All User's List
// @Description		Returns All Users in current user's organization,omits security sensitive fields
// @ID				GetAllUsers
// @Consume 		application/json
// @Produce 		application/json
// @Param			Cookie header string false "Send tokens using Cookie header<br>(Example:Cookie: accesstoken=eyJhbGc..; refreshtoken=eyJhbGciOiJ..)"
// @Tags			Access Users inventory
// @Security 		ApiToken
// @Success			200 {array} models.User
// @Failure			401
// @Failure			500
// @Router 			/users [get]
func GetAllUsers() gin.HandlerFunc {
	//Returns all users belonging to requested user's organization
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		//Filter by organization,all documents will be queried where user.org== docs.org
		filter := bson.D{{Key: "org", Value: ctx.GetString("org")}}
		//exclude password field
		opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}, {Key: "jwttoken", Value: 0}, {Key: "refreshtoken", Value: 0}})
		cursor, err := userCollection.Find(c, filter, opts)
		defer cancel()
		if err != nil {
			log.Panic(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var results []models.User //storing result in user list
		defer cancel()
		if err = cursor.All(c, &results); err != nil {
			log.Panic()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, results)

		//PAGINATION of records
		// 	recordsPerPage, err := strconv.Atoi(ctx.Query("records-per-page"))
		// 	if err != nil || recordsPerPage < 1 {
		// 		recordsPerPage = 10 // a default of 10 records per page
		// 	}
		// 	page, pageErr := strconv.Atoi(ctx.Query("page"))
		// 	if pageErr != nil || page < 1 {
		// 		page = 1
		// 	}

		// 	startIndex := (page - 1) * recordsPerPage
		// 	startIndex, err = strconv.Atoi(ctx.Query("start-index"))

		// 	//Pipeline functions
		// 	matchStage := bson.D{{"$match", bson.D{{}}}}
		// 	//used like filter
		// 	groupStage := bson.D{{"$group", bson.D{
		// 		{"_id", bson.D{{"_id", "null"}}},
		// 		{"totalcount", bson.D{{"$sum", 1}}},
		// 		{"$data", bson.D{{"$push", "$$ROOT"}}},
		// 	}}}
		// 	projectStage := bson.D{
		// 		{"$project", bson.D{
		// 			{"_id", 0},
		// 			{"totalcount", 1},
		// 			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordsPerPage}}}},
		// 		}},
		// 	}
		// 	result, err := userCollection.Aggregate(c, mongo.Pipeline{
		// 		matchStage,
		// 		groupStage,
		// 		projectStage,
		// 	})
		// 	defer cancel()
		// 	if err != nil {
		// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	}

		// 	var allUsersList []bson.M
		// 	if err = result.All(c, &allUsersList); err != nil {
		// 		log.Fatal(err)
		// 		return
		// 	}
		// 	ctx.JSON(http.StatusOK, allUsersList)

	}
}
