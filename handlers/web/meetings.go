package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

func Meetings() gin.HandlerFunc {
	return func(c *gin.Context) {
		classID := c.Param("classID")

		var jadwal []models.Jadwal
		if err := config.DB.
			Preload("Ruangan").
			Preload("Sesi").
			Preload("Pertemuan").
			Where("id_kelas = ?", classID).Find(&jadwal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Status": "Berhasil",
			"data": jadwal})

	}
}
