package middlewares

import (
	"context"
	"net/http"
	"time"

	AppConstant "github.com/HousewareHQ/backend-engineering-octernship/api/server/constants"
	"github.com/HousewareHQ/backend-engineering-octernship/api/server/helpers"
	"github.com/gin-gonic/gin"
)

// Middleware:Authenticating provided request's session
func SessionAuthentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		storedToken, storedTokenErr := ctx.Cookie("accesstoken") //get access-token cookie
		defer cancel()
		if storedTokenErr != nil {
			//if fails to get token from cookie,get it from header:Authorization
			storedToken = ctx.Request.Header.Get("Authorization") //TO GET TOKEN USING HEADER

		}
		//TRY:Validating access-token
		_, validTokenErr := helpers.ValidateJWTToken(storedToken)
		defer cancel()
		if validTokenErr != nil {
			if validTokenErr.Error() == "user no longer exists" {
				//If user doesnt exists anymore then user is unauthorized
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": validTokenErr.Error()})
				ctx.Abort()
				return
			}
			//accesstoken invalid or expired
			//Try to use refresh token to issue new access token
			storedRefreshToken, refreshTokenErr := ctx.Cookie("refreshtoken")
			defer cancel()

			if refreshTokenErr != nil {
				//IF:refresh token is not stored in cookies
				//THEN:Try to get it from header:Authorization
				storedRefreshToken = ctx.Request.Header.Get("Authorization") //TO GET TOKEN USING HEADER

				if storedRefreshToken == "" { //If that also fails then suggest relogin
					ctx.JSON(http.StatusForbidden, gin.H{"SessionExpired": refreshTokenErr.Error()})
					ctx.Abort()
					return
				}

			}
			//TRY:Validating refresh token
			_, validRefreshTokenErr := helpers.ValidateJWTToken(storedRefreshToken)
			defer cancel()
			if validRefreshTokenErr != nil {
				//IF:refreshtoken ,invalid or expired
				//REQUIRED: user to relogin ,SESSION EXPIRED!
				ctx.JSON(http.StatusForbidden, gin.H{"SessionExpired": validRefreshTokenErr})
				ctx.Abort()
				return
			}

			//If refreshToken is valid, DO:Refresh token rotation
			/*Issue a new access token and refresh token */
			accessToken, refreshToken, generateTokenErr := helpers.GenerateTokenByRefreshToken(c, storedRefreshToken)
			defer cancel()
			if generateTokenErr != nil {
				//IF:fail to generate valid token then user is no longer authorized.
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": generateTokenErr.Error()})
				ctx.Abort()
				return
			}

			//Setting tokens in cookies and in context
			ctx.SetCookie("accesstoken", accessToken, int(AppConstant.TOKEN_COOKIE_EXPIRY), "/", "localhost", false, true)
			ctx.SetCookie("refreshtoken", refreshToken, int(AppConstant.REFRESH_TOKEN_COOKIE_EXPIRY), "/", "localhost", false, true)
			ctx.Set("accesstoken", accessToken)
			ctx.Set("refreshtoken", refreshToken)
			ctx.Next() //passing current context to next one

		}

	}
}

// Middleware:Authorization
func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Either get accesstoken from previous context or current context's stored cookies*/
		/*
			If session is refreshed,then this middleware is also triggered
			but as previous context does not contain new stored token cookies information
			Hence,we will store new tokens in context using .Set() and then pass this context to next context
			and that will be this middleware's context.
			At the end we will check whether context.GetString("accesstoken") contains access token or not
			if yes: then that wil be our token
			else : this is new request and session is valid hence cookie contains valid tokens to be accessed and to proceed
		*/
		var clientJwtToken string
		var tokenAccessErr error
		if clientJwtToken = ctx.GetString("accesstoken"); clientJwtToken == "" {
			clientJwtToken, tokenAccessErr = ctx.Cookie("accesstoken")
			if tokenAccessErr != nil { //If Fails to get access token from cookie
				//Get token from header:Authorization
				clientJwtToken = ctx.Request.Header.Get("Authorization") //TO GET TOKEN USING HEADER

			}
			if clientJwtToken == "" { //if fails then suggest relogin

				ctx.JSON(http.StatusForbidden, gin.H{"SessionExpired": tokenAccessErr.Error()})
				ctx.Abort()
				return
			}

		}

		//Validate JWT token
		claims, err := helpers.ValidateJWTToken(clientJwtToken)
		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"TokenError": err})
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
