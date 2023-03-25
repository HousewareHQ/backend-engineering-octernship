package tokens

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	OrgId    uint   `json:"org_id"`
	jwt.StandardClaims
}
