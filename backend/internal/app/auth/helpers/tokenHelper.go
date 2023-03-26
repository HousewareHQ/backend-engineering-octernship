package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/database"

	jwt "github.com/dgrijalva/jwt-go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SignedDetails is the struct that will be encoded to a JWT.
type SignedDetails struct {
	Username string
	Uid string
	Orgid string
	UserType string
	jwt.StandardClaims
}
// create a new mongodb client
var userTokenCollection *mongo.Collection = database.OpenCollection(database.Client, "usertokens")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

// function input username, uid, orgid, usertype and return signedToken, signedRefreshToken, err
func GenerateAllTokens(username string, uid string, orgid string, usertype string)(signedToken string, signedRefreshToken string, err error){
	claims := &SignedDetails{
		Username: username,
		Uid: uid,
		Orgid: orgid,
		UserType: usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour).Unix(),// 1 hour
		},
	} 
	refreshClaims := &SignedDetails{
		Username: username,
		Uid: uid,
		Orgid: orgid,
		UserType: usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),// 24 hours
		},
	}
	// generate token and refresh token 
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	
	if err != nil {
		log.Panic(err)
		return 
	}
	return token, refreshToken, err
}

// ValidateToken is a function that takes a signed token string as input and returns a pointer to SignedDetails and a message string as output.
func ValidateToken(signedToken string) (claims *SignedDetails, msg string){
	// It parses the signed token and verifies it using the secret key. If the token is invalid, it returns an error message.
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{}, 
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)// claims: a pointer to SignedDetails that represents the claims of the signed token.
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		return
	}
	// If the token is expired, it returns an error message indicating that the token is expired.
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("the token is expired")
		return
	}
	return claims, msg
}


// UpdateAllTokens updates the user's tokens in the database.
// It takes as input the signedToken, signedRefreshToken, and userId.
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	// create a context with a 100 second timeout
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// create a primitive.D object to store the fields to be updated in the database
	var updateObj primitive.D

	// add the signedToken and signedRefreshToken fields to the update object
	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refreshtoken", signedRefreshToken})

	// get the current time and add it to the update object
	UpdatedAt, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
		return
	}
	updateObj = append(updateObj, bson.E{"updatedat", UpdatedAt})

	// set the options for the update operation, including the upsert option to insert a new document if the filter doesn't match any existing document
	upsert := true
	filter := bson.M{"userid":userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	// update the userTokenCollection with the new tokens and update time
	_, err = userTokenCollection.UpdateMany(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	// close the context to prevent memory leaks
	defer cancel()

	// check for errors and log them
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}
