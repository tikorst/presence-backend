package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		jwtClaims := claims.(jwt.MapClaims)
		role := jwtClaims["role"].(string)

		c.JSON(http.StatusOK, gin.H{"role": role})
	}
}
