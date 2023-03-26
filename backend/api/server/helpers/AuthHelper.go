package helpers

import (
	"github.com/gin-gonic/gin"
)

// Returns true if user is admin
func IsAdmin(ctx *gin.Context) bool {
	return ctx.GetString("usertype") == "ADMIN"
}

// Return true if current admin and user getting deleted user belongs to same organization
func AreFromSameOrg(userOrg string, deleteUserOrg string) bool {
	return userOrg == deleteUserOrg
}
