package mobile

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}

