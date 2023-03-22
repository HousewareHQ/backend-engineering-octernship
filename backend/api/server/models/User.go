package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model
type User struct {
	Username     string             `json:"username" validate:"required,min=3,max=15"`
	Password     string             `json:"password" validate:"required,min=4"`
	Usertype     string             `json:"usertype" validate:"required,eq=USER|eq=ADMIN"`
	Org          string             `json:"org"`
	JWTToken     string             `json:"jwttoken"`
	RefreshToken string             `json:"refreshtoken"`
	CreatedOn    time.Time          `json:"createdon"`
	UpdatedOn    time.Time          `json:"updatedon"`
	ID           primitive.ObjectID `bson:"_id"`
}
