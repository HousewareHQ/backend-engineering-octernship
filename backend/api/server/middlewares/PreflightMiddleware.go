package middlewares

import "github.com/gin-gonic/gin"

func PreflightRequestMiddleware() gin.HandlerFunc { //preflight request resource policy
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") //ORIGIN:From where API calls will be made i.e.React.js client.
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")             //TO TRANSFER COOKIES
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) //No content success:Client doesnt need to navigate away from current state
			return
		}

		c.Next() //pass context to next one
	}
}
