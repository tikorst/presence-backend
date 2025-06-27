package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/helpers"
)

// VerifyRole handles the request to verify the user's role
func VerifyRole(c *gin.Context) {

	// Get the user's role from the context
	role, _ := helpers.GetRole(c)

	// return the role in the response
	c.JSON(http.StatusOK, gin.H{"role": role})
}
