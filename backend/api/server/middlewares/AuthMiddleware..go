package middlewares

import (
	"net/http"

	"github.com/HousewareHQ/backend-engineering-octernship/api/server/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientJwtToken := ctx.Request.Header.Get("token")

		if clientJwtToken == "" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorizd:User not logged in"})
			ctx.Abort() // Abort all operations related to this context.
			return
		}
		claims, err := helpers.ValidateJWTToken(clientJwtToken)
		if err != "" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err})
			ctx.Abort()
			return
		}
		//Setting context's data with claims,can be used by functions accessing background contexts
		ctx.Set("username", claims.Username)
		ctx.Set("usertype", claims.Usertype)
		ctx.Set("createdon", claims.CreatedAt)

		ctx.Next() //passing context to next handler
	}
}
