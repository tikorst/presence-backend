package mobile

import (
	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

// GetSemester handles the request to get all semesters
func GetSemester(c *gin.Context) {

	// Query to get all semesters from the database
	var semesters []models.Semester
	if err := config.DB.Find(&semesters).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data semester"})
		return
	}

	// Return the list of semesters
	c.JSON(200, gin.H{"erorr": false, "data": semesters})

}
