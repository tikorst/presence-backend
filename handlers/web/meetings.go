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
			Preload("Kelas").
			Preload("Kelas.MataKuliah").
			Preload("Pertemuan").
			Where("id_kelas = ?", classID).Find(&jadwal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
			return
		}

		// Extract all IDJadwal from the jadwal slice
		// var jadwalIDs []int
		// for _, j := range jadwal {
		// 	jadwalIDs = append(jadwalIDs, j.IDJadwal)
		// }

		// // Fetch all pertemuan based on the extracted IDJadwal
		// var pertemuan []models.Pertemuan
		// if err := config.DB.Debug().Where("id_jadwal = ?", jadwal[0].IDJadwal).
		// 	Preload("Jadwal.Ruangan").
		// 	Preload("Jadwal.Sesi").
		// 	Preload("Jadwal.Kelas").Find(&pertemuan).Error; err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meetings"})
		// 	return
		// }

		c.JSON(http.StatusOK, gin.H{"Status": "Berhasil",
			"data": jadwal})

	}
}
