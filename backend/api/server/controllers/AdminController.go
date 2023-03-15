package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	helpers "github.com/HousewareHQ/backend-engineering-octernship/api/server/helpers"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// CREATE USER (ONLY ADMIN)
func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Cancel context,after 10 sec (TIME-OUT)
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//ONLY ADMIN CAN CREATE A ACCOUNT
		if helpers.IsAdmin(ctx) {
			var newUser models.User

			if err := ctx.BindJSON(&newUser); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

				return
			}

			//DBconnect.PingDB()->To test connection to DB

			//Validating:Does user exists already?
			count, err := userCollection.CountDocuments(context.TODO(), bson.M{"username": newUser.Username})
			if err != nil {

				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if count > 1 {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists"})
				return
			}

			//If does not exists,then create new user
			newUser.CreatedOn, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			newUser.UpdatedOn = newUser.CreatedOn
			token, refreshToken := helpers.GenerateTokens(newUser)
			newUser.JWTToken = token
			newUser.RefreshToken = refreshToken
			newUser.ID = primitive.NewObjectID()
			newPass := newUser.Password
			newUser.Password = helpers.Hashing(newPass) //hashing password

			if err := validate.Struct(newUser); err != nil { //Validating user struct
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			resInsNum, insertErr := userCollection.InsertOne(c, newUser) //storing in db
			if insertErr != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"result_insertion_number": resInsNum})

		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized logon"})
			return
		}
	}

}

/*DELETE A USER */
func DeleteUser() gin.HandlerFunc {
	//ONLY ADMIN CAN MODIFY/DELETE USER

	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		uid := ctx.Param("uid") //getting parameter uid from url
		userObjectId, err := primitive.ObjectIDFromHex(uid)
		defer cancel()

		if err != nil {
			log.Panic("Incorrect User ObjectID:MongoDB")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		filter := bson.D{{Key: "_id", Value: userObjectId}}
		res, delErr := userCollection.DeleteOne(c, filter)
		defer cancel()
		if delErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": delErr.Error()})
			return
		}
		ctx.JSON(http.StatusOK, res)
	}

}
