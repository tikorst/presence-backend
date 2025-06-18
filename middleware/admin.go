package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/helpers"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := helpers.GetRole(c)
		if role != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access restricted to Admins only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
