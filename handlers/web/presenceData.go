package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

// func GetPresenceData retrieves the attendance data for a specific meeting and class.
func GetPresenceData(c *gin.Context) {

	// Get the meetingID and classID from the URL parameters
	meetingID := c.Param("meetingID")
	classID := c.Param("classID")

	// Query to get all students registered in the class
	var mahasiswaKelas []models.MahasiswaKelas
	if err := config.DB.Preload("Mahasiswa").Preload("Mahasiswa.User").Where("id_kelas = ? AND status = ?", classID, "aktif").Find(&mahasiswaKelas).Order("npm DESC").Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak ada mahasiswa yang terdaftar di kelas ini"})
		return
	}

	// Query to get all attendance records for the specified meeting
	var presensi []models.Presensi
	if err := config.DB.Where("id_pertemuan = ?", meetingID).Find(&presensi).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Belum ada presensi untuk pertemuan ini"})
		return
	}

	// Create maps to hold attendance status and notes
	attendanceMap := make(map[string]string)
	catatanMap := make(map[string]string)

	// Initialize attendance status to "Alpha" and notes to empty for all students
	for _, mhs := range mahasiswaKelas {
		attendanceMap[mhs.Mahasiswa.User.Nama] = "Alpha"
		catatanMap[mhs.NPM] = ""
	}

	// Update the attendance status and notes for students who are present
	// based on the attendance records
	for _, pres := range presensi {
		attendanceMap[pres.NPM] = "Hadir"
		catatanMap[pres.NPM] = pres.Catatan
	}

	// Create a response data structure
	responseData := make([]map[string]string, 0)
	for _, mhs := range mahasiswaKelas {
		responseData = append(responseData, map[string]string{
			"npm":    mhs.NPM,
			"nama":   mhs.Mahasiswa.User.Nama,
			"status": attendanceMap[mhs.NPM],
			"catatan": catatanMap[mhs.NPM], 
		})
	}


	// Return the attendance data as JSON
	c.JSON(http.StatusOK, gin.H{"attendance": responseData})
}
