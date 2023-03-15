package controllers

import (
	"context"
	"net/http"
	"time"

	AppConstant "github.com/HousewareHQ/backend-engineering-octernship/api/server/constants"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/helpers"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/models"
	DBconnect "github.com/HousewareHQ/backend-engineering-octernship/api/server/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = DBconnect.OpenCollection(DBconnect.Client, AppConstant.USER_COLLECTION)
var validate = validator.New()

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var storedUser models.User
		defer cancel()
		//unmarshal encoded-json into struct
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		helpers.UpdateTokenOnLogin(token, refreshedToken, storedUser.ID)

		ctx.JSON(http.StatusOK, storedUser)

	}

}

func GetAllUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		// //PAGINATION of records
		// recordsPerPage, err := strconv.Atoi(ctx.Query("records-per-page"))
		// if err != nil || recordsPerPage < 1 {
		// 	recordsPerPage = 10 // a default of 10 records per page
		// }
		// page,pageErr := strconv.Atoi(ctx.Query("page"))
		// if pageErr != nil || page <1{
		// 	page =1
		// }

		// startIndex:= (page-1)*recordsPerPage
		// startIndex,err :=

	}
}
