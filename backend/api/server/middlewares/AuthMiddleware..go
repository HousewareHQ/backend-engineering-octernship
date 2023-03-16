package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	AppConstant "github.com/HousewareHQ/backend-engineering-octernship/api/server/constants"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/helpers"
	"github.com/gin-gonic/gin"
)

func SessionAuthentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		storedToken, tokenerr := ctx.Cookie("accesstoken")
		defer cancel()
		if tokenerr != nil {
			//access token is not stored in cookie
			//assuming refreshtoken is also not stored
			//Suggesting client to relogin
			ctx.JSON(http.StatusForbidden, gin.H{"SessionExpired": tokenerr.Error()})
			ctx.Abort()
			return

		}
		_, validTokenErr := helpers.ValidateJWTToken(storedToken)
		if validTokenErr != "" {
			//accesstokeninvalid or expired
			//Try to use refresh token to issue new access token
			storedRefreshToken, refreshTokenErr := ctx.Cookie("refreshtoken")
			defer cancel()
			if refreshTokenErr != nil {
				//refresh token is not stored
				//Suggesting client to relogin
				ctx.JSON(http.StatusForbidden, gin.H{"SessionExpired": refreshTokenErr.Error()})
				ctx.Abort()
				return

			}
			_, validRefreshTokenErr := helpers.ValidateJWTToken(storedRefreshToken)
			defer cancel()
			if validRefreshTokenErr != "" {
				//refreshtoken ,invalid or expired
				//Suggest user to login again,SESSION EXPIRED!
				ctx.JSON(http.StatusForbidden, gin.H{"SessionExpired": validRefreshTokenErr})
				ctx.Abort()
				return
			}
			//If refreshToken is valid
			/*Issue a new access token and refresh token */
			accessToken, refreshToken, generateTokenErr := helpers.GenerateTokenByRefreshToken(c, storedRefreshToken)
			defer cancel()
			if generateTokenErr != nil {
				log.Panic("Error while generating token using refresh token")
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": generateTokenErr.Error()})
				ctx.Abort()
				return
			}
			ctx.SetCookie("accesstoken", accessToken, int(AppConstant.TOKEN_EXPIRY), "/users", "localhost", false, true)
			ctx.SetCookie("refreshtoken", refreshToken, int(AppConstant.REFRESH_TOKEN_EXPIRY), "/users", "localhost", false, true)
			ctx.JSON(http.StatusAccepted, gin.H{"Ok": "Refreshed session"})

		}

	}
}

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// clientJwtToken := ctx.Request.Header.Get("token")
		clientJwtToken, tokenAccessErr := ctx.Cookie("accesstoken")
		if tokenAccessErr != nil { //If Fails to get access token from cookie:Suggest relogin
			ctx.JSON(http.StatusForbidden, gin.H{"SessionExpired": tokenAccessErr.Error()})
			ctx.Abort()
			return
		}

		//Validate JWT token
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
		ctx.Set("org", claims.Org)

		ctx.Next() //passing context to next handler
	}
}
