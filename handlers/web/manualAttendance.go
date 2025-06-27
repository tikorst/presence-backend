package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
)

type ManualAttendanceRequest struct {
	NPM     string `json:"npm" binding:"required"`
	Catatan string `json:"catatan"`
}

func ManualAttendance(c *gin.Context) {

	// get classID and meetingID from URL parameters
	classIDStr := c.Param("classID")
	meetingIDStr := c.Param("meetingID")

	// Get the username from the context
	username, _ := helpers.GetUsername(c)

	// Convert classID and meetingID to integers
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
		return
	}

	// Convert meetingID to integer
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	// Check if the user is a lecturer for the class
	var dosenPengampu models.DosenPengampu
	if err := config.DB.Where("id_kelas = ? AND nip = ?", classID, username).First(&dosenPengampu).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk kelas ini"})
		return
	}

	// Bind request data
	var req ManualAttendanceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Verify if the student is registered in the class
	var mahasiswaKelas models.MahasiswaKelas
	if err := config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", req.NPM, classID, "aktif").First(&mahasiswaKelas).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mahasiswa tidak terdaftar di kelas ini"})
		return
	}

	// Check if the student has already checked in for this meeting
	var existingPresensi models.Presensi
	if err := config.DB.Where("npm = ? AND id_pertemuan = ?", req.NPM, meetingID).First(&existingPresensi).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Mahasiswa sudah melakukan presensi untuk pertemuan ini"})
		return
	}

	// Insert attendance record
	presensi := models.Presensi{
		NPM:           req.NPM,
		IDPertemuan:   meetingID,
		WaktuPresensi: time.Now(),
		Status:        "Hadir",
		Catatan:       req.Catatan,
	}

	// Set the class ID for the attendance record
	if err := config.DB.Create(&presensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data presensi"})
		return
	}

	// return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Presensi berhasil ditambahkan",
		"data":    presensi,
	})
}
