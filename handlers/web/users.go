package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		jwtClaims := claims.(jwt.MapClaims)
		username := jwtClaims["sub"].(string)

		currentUser := models.User{}

		if err := config.DB.Debug().Where("username = ?", username).First(&currentUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch current user"})
			return
		}

		if currentUser.TipeUser != "Admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		users := []models.User{}

		if err := config.DB.Debug().Where("tipe_user = ?", "Mahasiswa").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
