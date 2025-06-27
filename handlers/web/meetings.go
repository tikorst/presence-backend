package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

func GetMeetings(c *gin.Context) {

	// Get the classID from the URL parameters
	classID := c.Param("classID")

	// Query to get the schedules for the class
	var jadwal []models.Jadwal
	if err := config.DB.
		Preload("Ruangan").
		Preload("Sesi").
		Preload("Pertemuan").
		Where("id_kelas = ?", classID).Find(&jadwal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}

	// return the schedules
	c.JSON(http.StatusOK, gin.H{"Status": "Berhasil",
		"data": jadwal})

}

