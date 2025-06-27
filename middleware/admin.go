package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/helpers"
)

// AdminOnly is a middleware function that restricts access to routes for Admin users only.
func AdminOnly() gin.HandlerFunc {

	// Return a middleware function that checks the user's role
	return func(c *gin.Context) {
		role, _ := helpers.GetRole(c)
		if role != "Admin" {

			// If the role is not Admin, return a 403 Forbidden response
			c.JSON(http.StatusForbidden, gin.H{"error": "Access restricted to Admins only"})
			c.Abort()
			return
		}

		// If the role is Admin, proceed to the next handler
		c.Next()
	}
}
