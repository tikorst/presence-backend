package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AppSecretMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.GetHeader("X-App-Secret")
		expected := os.Getenv("APP_SECRET_KEY")

		if secret != expected {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Unauthorized app access",
			})
			return
		}
		fmt.Printf("App Secret verified: %s\n", secret)
		c.Next()
	}
}
