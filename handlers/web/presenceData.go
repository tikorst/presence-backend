package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

func GetPresenceData(c *gin.Context)   {
	meetingID := c.Param("meetingID")
	classID := c.Param("classID")
	var mahasiswaKelas []models.MahasiswaKelas
	if err := config.DB.Preload("Mahasiswa").Preload("Mahasiswa.User").Where("id_kelas = ? AND status = ?", classID, "aktif").Find(&mahasiswaKelas).Order("npm DESC").Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak ada mahasiswa yang terdaftar di kelas ini"})
		return
	}
	var presensi []models.Presensi
	if err := config.DB.Where("id_pertemuan = ?", meetingID).Find(&presensi).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Belum ada presensi untuk pertemuan ini"})
		return
	}
	attendanceMap := make(map[string]string)
	for _, mhs := range mahasiswaKelas {
		attendanceMap[mhs.Mahasiswa.User.Nama] = "Alpha"
	}

	for _, pres := range presensi {
		attendanceMap[pres.NPM] = "Hadir"
	}

	// Create a response data structure
	responseData := make([]map[string]string, 0)
	for _, mhs := range mahasiswaKelas {
		responseData = append(responseData, map[string]string{
			"nama":   mhs.Mahasiswa.User.Nama,
			"status": attendanceMap[mhs.NPM],
		})
	}

	c.JSON(http.StatusOK, gin.H{"attendance": responseData})
}
