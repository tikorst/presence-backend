package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
)

// GetUsers handles the request to fetch users with pagination and search functionality
func GetUsers(c *gin.Context) {

	// Get the username from the context
	username, _ := helpers.GetUsername(c)

	// Initialize the currentUser variable to check if the user is an admin
	currentUser := models.User{}

	// Query to get the current user based on the username
	if err := config.DB.Where("username = ?", username).First(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch current user"})
		return
	}

	// Check if the current user is an admin
	// If not, return an unauthorized error
	if currentUser.TipeUser != "Admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get pagination and search parameters from the query string
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	searchQuery := c.DefaultQuery("search", "")

	// Convert page and limit to integers, with default values if conversion fails
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Convert limit to integer, with a default value if conversion fails
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Calculate the offset for pagination
	offset := (page - 1) * limit

	// Initialize variables to hold the users and total count
	var users []models.User
	var totalUsers int64

	// Create a query to fetch users with the specified type
	dbQuery := config.DB.Where("tipe_user = ?", "Mahasiswa")

	// If a search query is provided, add a condition to filter users by username or name
	if searchQuery != "" {
		dbQuery = dbQuery.Where("username LIKE ? OR nama LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Count the total number of users matching the query
	if err := dbQuery.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		fmt.Println("Error counting users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

	// Fetch the users with pagination and search applied
	if err := dbQuery.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		fmt.Println("Error fetching users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Calculate the total number of pages based on the total users and limit
	totalPages := (totalUsers + int64(limit) - 1) / int64(limit)

	// Return response with users, totalUsers, totalPages, currentPage, and limit
	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"totalUsers":  totalUsers,
		"totalPages":  totalPages,
		"currentPage": page,
		"limit":       limit,
	})
}
