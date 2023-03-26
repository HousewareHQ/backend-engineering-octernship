package middleware

import (
	"fmt"
	"net/http"

	helper "github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/helpers"

	"github.com/gin-gonic/gin"
)

// Authenticate is a middleware function that verifies if the client is authorized to access a specific resource.
// It returns a gin.HandlerFunc.
func Authenticate() gin.HandlerFunc{
	return func(c *gin.Context){
		// Retrieve the token from the Authorization header of the request.
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return 
		}
		// Validate the token using the ValidateToken function from the helper package.
		claims, err := helper.ValidateToken(clientToken)
		if err != ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return 
		}
		if claims.UserType != "ADMIN"{
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Unauthorized to access this resource")})
			c.Abort()
			return 
		}
		// Set the "username" and "uid" values in the Gin context for use in subsequent requests.
		c.Set("username", claims.Username)
		c.Set("uid", claims.Uid)
		c.Next()// Call the next middleware function in the chain.
	}
}

// TokenAuth is a middleware function that validates the token received in the Authorization header of the incoming request
func TokenAuth() gin.HandlerFunc{
	return func(c *gin.Context){
		clientToken := c.Request.Header.Get("token")// get the client token from the Authorization header
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return 
		}
		// validate the client token
		claims, err := helper.ValidateToken(clientToken)
		if err != ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return 
		}
		// match the client token with the user token in the database
		errr := helper.MatchClientTokenToUserToken(claims.Uid, clientToken)
		if errr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user session logged out, you need to login"})
			c.Abort()
			return
		}
		c.Set("orgId", claims.Orgid)
		c.Next()
	}
}