package mobile

import (
	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

func Semester(c *gin.Context) {

	var semesters []models.Semester
	if err := config.DB.Find(&semesters).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data semester"})
		return
	}
	c.JSON(200, gin.H{"erorr": false, "data": semesters})

}
