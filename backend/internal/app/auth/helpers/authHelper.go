package helpers

import (
	"context"
	"errors"

	models "github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/models"
	"go.mongodb.org/mongo-driver/bson"
)

// MatchClientTokenToUserToken function compares the client token with the user token stored in the database for a given user ID.
func MatchClientTokenToUserToken(userid string, clientToken string) error {
	var user models.UserTokens
	filter := bson.M{"userid": userid}
	err := userTokenCollection.FindOne(context.Background(), filter).Decode(&user)
	// if the token is not found in the database, return an error
	if err != nil {
		return err
	}
	// The function returns an error if the user session is logged out, otherwise it returns nil.

	if *user.Token != clientToken{
		return errors.New("user session logged out, you need to login")
	}
	return nil
}
