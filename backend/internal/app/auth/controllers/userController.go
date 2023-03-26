package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/database"
	helper "github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/helpers"
	"github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var userTokenCollection *mongo.Collection = database.OpenCollection(database.Client, "usertokens") 
var validate = validator.New() // validator for the user model


func isValidPassword(password string) bool {
    if len(password) < 8 {
        return false
    }
    hasUpperCase := false
    hasLowerCase := false
    hasDigit := false
    hasSpecialChar := false
    for _, ch := range password {
        switch {
        case 'A' <= ch && ch <= 'Z':
            hasUpperCase = true
        case 'a' <= ch && ch <= 'z':
            hasLowerCase = true
        case '0' <= ch && ch <= '9':
            hasDigit = true
        case strings.ContainsAny(string(ch), "!@#$%^&*"):
            hasSpecialChar = true
        }
    }
    return hasUpperCase && hasLowerCase && hasDigit && hasSpecialChar
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // generate a hashed password
	if err != nil {
		log.Panic(err)
	}
	return string(bytes), err
}

func VerifyPassword(userPassword string, providedPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword)) // compare the hashed password with the provided password
	check := true
	msg := ""
	
	if err!= nil {
		msg = fmt.Sprintf("Username or password is incorrect")
		check=false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var usertokens models.UserTokens

		
		if err := c.BindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}
		
		// below code is to validate the user variable			
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// hash the password and unit testing 
		password, err := HashPassword(*user.Password) 
		user.Password = &password
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while hashing the password"})
		}

		// if username is already present then print error "username already exists"
		count, err := userCollection.CountDocuments(ctx, bson.M{"username":user.Username})
		defer cancel()
		if err!= nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for the username"})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		user.CreatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		usertokens.CreatedAt = user.CreatedAt
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing the time"})
		}

		user.UpdatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		usertokens.UpdatedAt = user.UpdatedAt
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing the time"})
		}
		user.ID = primitive.NewObjectID()
		usertokens.ID = primitive.NewObjectID()

		user.UserID = user.ID.Hex()
		usertokens.UserID = user.UserID // mapping the user_id of the user to the user_id of the usertokens

		user.UserType = "ADMIN"

		token, refreshToken, err := helper.GenerateAllTokens(*user.Username, *&user.UserID, *&user.OrgID, *&user.UserType)
		if err != nil {
			log.Panic(err) 
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while generating the tokens"})
		}
		usertokens.Token = &token
		usertokens.RefreshToken = &refreshToken

		resultInsertionNumberUser, insertErrUser := userCollection.InsertOne(ctx, user) // insert the user into the database
		resultInsertionNumberUserTokens, insertErrUserTokens := userTokenCollection.InsertOne(ctx, usertokens) // insert the usertokens into the database
		if insertErrUser != nil {
			msg := fmt.Sprintf("user item not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if insertErrUserTokens != nil {
			msg := fmt.Sprintf("usertokens item not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"userInsertionNumber": resultInsertionNumberUser,
			"userTokensInsertionNumber": resultInsertionNumberUserTokens,
		})
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) 
		var user models.User
		var foundUser models.User
		var foundUserToken models.UserTokens
		  
		//bind the request body to the user struct
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//check if the user exists in the database
		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		defer cancel()
		if foundUser.Username == nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "username not found"})
			return
		}
		// if the user is not found, return an error
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Username or password is incorrect"})
			return
		}
		// verify the password by calling the verify password function
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		//if the password is not valid, return an error
		if passwordIsValid == false{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		// generate new tokens when logging in
		token, refreshToken, errr := helper.GenerateAllTokens(*foundUser.Username, foundUser.UserID, foundUser.OrgID, foundUser.UserType)
		if errr != nil {
			log.Panic(errr) 
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while generating the tokens"})
		}
		// update the tokens in the database
		helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		// return the user and the tokens
		err = userCollection.FindOne(ctx, bson.M{"userid":foundUser.UserID}).Decode(&foundUser)
		err = userTokenCollection.FindOne(ctx, bson.M{"userid":foundUser.UserID}).Decode(&foundUserToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"user": foundUser,
			"Access Token": token,
			"Refresh Token": refreshToken,
		})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) // create a context with a timeout

		orgId := c.GetString("orgId") // retrieve the orgId from the context

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage")) // get the limit from the url
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10 // set default limit to 10 if it is not provided or invalid
		}

		page, err := strconv.Atoi(c.Query("page")) // get the page from the url
		if err != nil || page < 1 {
			page = 1 // set default page to 1 if it is not provided or invalid
		}

		startIndex, err := strconv.Atoi(c.Query("startIndex")) // get the start index from the url
		if err != nil || startIndex < 0 {
			startIndex = 0 // set default start index to 0 if it is not provided or invalid
		}

		// aggregation pipeline to filter users by orgId and paginate the results
		matchStage := bson.D{{"$match", bson.M{"orgid": orgId}}} //$match stage filters the documents based on the specified orgId value
		groupStage := bson.D{{"$group", bson.D{ //$group stage, which groups the filtered documents by dummy _id field 
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}}, // sum to count totoal number of documents
			{"data", bson.D{{"$push", "$$ROOT"}}}}}} // push all the documents to data field
		projectStage := bson.D{ // modifies the output of the aggregation pipeline
			{"$project", bson.D{
				{"_id", 0}, // exclude _id field from the output
				{"total_count", 1},	// include total_count field in the output
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}	// include user_items field in the output

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage}) // aggregate the data
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil { // decode the data
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allUsers) // return the paginated results
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the request header
		tokenString := c.GetHeader("Token")

		// Validate the token
		claims, msg := helper.ValidateToken(tokenString)
		if claims == nil {
			c.JSON(401, gin.H{"error": msg})
			return
		}
		// Update the user's tokens
		userId := claims.Uid
		helper.UpdateAllTokens("", "", userId)

		// Return success message
		c.JSON(200, gin.H{"message": "User logged out successfully"})
	}
}

func AddUserToOrg() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) // create a context with a timeout
		
		var user models.User // create a user variable of type User
		var usertokens models.UserTokens // create a usertokens variable of type UserTokens
		
		// below code is to extract the org_id from the token 
		clientToken:= c.Request.Header.Get("token") // get the token from the request header
		claims, msg := helper.ValidateToken(clientToken) // validate the token
		if claims == nil {
			c.JSON(401, gin.H{"error": msg})
			return
		}
		if msg != ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return 
		}
		user.OrgID = claims.Orgid // set the org_id of the user to the org_id in the token
		
		// below code is to bind extract the json payload from the request body and map it to the 'user' struct
		if err := c.BindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}
		// below code is to validate the user variable			
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		// below code is to check if the password contains atleast one alphabet, one number and one special character
		if !isValidPassword(*user.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is invalid, create a new password with at least one uppercase and lowercase alphabet, one number, and one special character from the set !@#$%^&* and minimum length of 8 characters"})
			return 
		}

		// hash the password and unit testing 
		password, err := HashPassword(*user.Password) 
		user.Password = &password
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while hashing the password"})
		}

		// if username is already present then print error "username already exists"
		count, err := userCollection.CountDocuments(ctx, bson.M{"username":user.Username})
		defer cancel()
		if err!= nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for the username"})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		user.CreatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		usertokens.CreatedAt = user.CreatedAt
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing the time"})
		}

		user.UpdatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		usertokens.UpdatedAt = user.UpdatedAt
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing the time"})
		}
		user.ID = primitive.NewObjectID()
		usertokens.ID = primitive.NewObjectID()

		user.UserID = user.ID.Hex()
		usertokens.UserID = user.UserID // mapping the user_id of the user to the user_id of the usertokens


		token, refreshToken, err := helper.GenerateAllTokens(*user.Username, *&user.UserID, *&user.OrgID, *&user.UserType)
		if err != nil {
			log.Panic(err) 
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while generating the tokens"})
		}
		usertokens.Token = &token
		usertokens.RefreshToken = &refreshToken

		resultInsertionNumberUser, insertErrUser := userCollection.InsertOne(ctx, user) // insert the user into the database
		resultInsertionNumberUserTokens, insertErrUserTokens := userTokenCollection.InsertOne(ctx, usertokens) // insert the usertokens into the database
		if insertErrUser != nil {
			msg := fmt.Sprintf("user item not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if insertErrUserTokens != nil {
			msg := fmt.Sprintf("usertokens item not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"userInsertionNumber": resultInsertionNumberUser,
			"userTokensInsertionNumber": resultInsertionNumberUserTokens,
		})
	}
}

func DeleteUserFromOrg() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// var user models.User
		var user models.User

		userID := c.Query("userId") // get the user_id from the query params
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not provided"})
			return
		}
		
		//below code is to extract the org_id from the token
		clientToken:= c.Request.Header.Get("token") // get the token from the request header
		claims, msg := helper.ValidateToken(clientToken) // validate the token
		if claims == nil {
			c.JSON(401, gin.H{"error": msg})
			return
		}
		if msg != ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return 
		}
		// below code is to find the user in the database which has the same org_id and user_id as the one in the url and in the token
		err := userCollection.FindOne(ctx, bson.M{"orgid": claims.Orgid, "userid": userID}).Decode(&user)
		// if the user is not found, return an error
		defer cancel()
		if user.Username == nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "wrong userID or the user doesn't belong in your organization"})
			return
		}
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}
		
		// delete the user from the user collection
		_, err = userCollection.DeleteOne(ctx, bson.M{"userid": userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}
		// delete from the usertokens collection
		_, errr := userTokenCollection.DeleteMany(ctx, bson.M{"userid": userID})
		if errr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func RefreshAccessToken() gin.HandlerFunc{
	return func(c *gin.Context){
		var user models.User
		refreshToken := c.Request.Header.Get("refreshtoken") // get the refresh token from the request header 
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token not provided"})
			return
		}

		claims, err := helper.ValidateToken(refreshToken) // validate the refresh token
		if err != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
			return
		}
		// check if refresh token is expired 
		if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token expired"})
			return
		}
		user.Username = &claims.Username
		user.UserID = claims.Uid
		user.OrgID = claims.Orgid
		user.UserType = claims.UserType

		NewAccessToken, _, errr := helper.GenerateAllTokens(*user.Username, *&user.UserID, *&user.OrgID, *&user.UserType) // generate a new access token
		if errr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}
		helper.UpdateAllTokens(NewAccessToken, refreshToken, user.UserID)
		c.JSON(http.StatusOK, gin.H{"accessToken": NewAccessToken})
	}
}