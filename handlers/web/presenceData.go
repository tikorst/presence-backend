package web

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

// func GetPresenceData retrieves the attendance data for a specific meeting and class.
func GetPresenceData(c *gin.Context) {

	// Get the meetingID and classID from the URL parameters
	meetingID := c.Param("meetingID")
	classID := c.Param("classID")

	// Initialize variables for concurrent queries
	// var mahasiswaKelas []models.MahasiswaKelas
	var presensi []models.Presensi
	var mahasiswaErr, presensiErr error
	var wg sync.WaitGroup
	type MahasiswaInfo struct {
	NPM  string
	Nama string
}

	var mahasiswaInfo []MahasiswaInfo
	// Goroutine for fetching students registered in the class
	wg.Add(2)
	go func() {
		defer wg.Done()
		mahasiswaErr = config.DB.
			Table("mahasiswa_kelas").
			Select("mahasiswa_kelas.npm, users.nama").
			Joins("JOIN mahasiswa ON mahasiswa.npm = mahasiswa_kelas.npm").
			Joins("JOIN users ON users.id_user = mahasiswa.id_user").
			Where("mahasiswa_kelas.id_kelas = ? AND mahasiswa_kelas.status = ?", classID, "aktif").
			Order("mahasiswa_kelas.npm DESC").
			Find(&mahasiswaInfo).Error
	}()

	// Goroutine for fetching attendance records for the meeting
	go func() {
		defer wg.Done()
		presensiErr = config.DB.
			Where("id_pertemuan = ?", meetingID).
			Find(&presensi).Error
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Check for errors
	if mahasiswaErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak ada mahasiswa yang terdaftar di kelas ini"})
		return
	}

	if presensiErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Belum ada presensi untuk pertemuan ini"})
		return
	}

	// Create maps to hold attendance status and notes
	attendanceMap := make(map[string]string)
	catatanMap := make(map[string]string)

	// Initialize attendance status to empty and notes to empty for all students
	for _, mhs := range mahasiswaInfo {
		attendanceMap[mhs.NPM] = ""
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
	for _, mhs := range mahasiswaInfo {
		responseData = append(responseData, map[string]string{
			"npm":    mhs.NPM,
			"nama":   mhs.Nama,
			"status": attendanceMap[mhs.NPM],
			"catatan": catatanMap[mhs.NPM], 
		})
	}

	// Return the attendance data as JSON
	c.JSON(http.StatusOK, gin.H{"attendance": responseData})
}