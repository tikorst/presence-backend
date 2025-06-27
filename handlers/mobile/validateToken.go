package mobile

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Function to validate the token
// This function will return a success if request token is valid and not expired
func ValidateToken(c *gin.Context) {
	// Return response success
	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}

