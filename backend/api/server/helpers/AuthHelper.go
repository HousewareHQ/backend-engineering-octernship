package helpers

import (
	"github.com/gin-gonic/gin"
)

func IsAdmin(ctx *gin.Context) bool {
	return ctx.GetString("usertype") == "ADMIN"
}
