package web

import (
	"fmt"
	"net/http"
	"strconv"

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

		if err := config.DB.Where("username = ?", username).First(&currentUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch current user"})
			return
		}

		if currentUser.TipeUser != "Admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")
		searchQuery := c.DefaultQuery("search", "")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		var users []models.User
		var totalUsers int64

		dbQuery := config.DB.Debug().Where("tipe_user = ?", "Mahasiswa")

		fmt.Printf("Backend received: Page=%d, Limit=%d, SearchQuery='%s'\n", page, limit, searchQuery)

		if searchQuery != "" {
			dbQuery = dbQuery.Where("username LIKE ? OR nama LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
		}

		if err := dbQuery.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
			fmt.Println("Error counting users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
			return
		}

		if err := dbQuery.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
			fmt.Println("Error fetching users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		totalPages := (totalUsers + int64(limit) - 1) / int64(limit)

		c.JSON(http.StatusOK, gin.H{
			"users":       users,
			"totalUsers":  totalUsers,
			"totalPages":  totalPages,
			"currentPage": page,
			"limit":       limit,
		})
	}
}