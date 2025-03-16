package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

func Classes() gin.HandlerFunc {
	return func(c *gin.Context) {
		var classes []models.DosenPengampu
		claims, _ := c.Get("claims")
		jwtClaims := claims.(jwt.MapClaims)
		username := jwtClaims["sub"].(string)

		if err := config.DB.Preload("Kelas.MataKuliah").Find(&classes).Where("nip = ?", username).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch classes"})
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"classes": classes})
	}
}
