package web

import (
	"net/http"
	"strconv"
	"sync"
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

	// Bind request data early
	var req ManualAttendanceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Variables for parallel queries
	var dosenPengampu models.DosenPengampu
	var mahasiswaKelas models.MahasiswaKelas
	var existingPresensi models.Presensi
	var dosenErr, mahasiswaErr, presensiErr error
	var wg sync.WaitGroup

	wg.Add(3)
	// Goroutine 1: Check if user is lecturer for the class
	go func() {
		defer wg.Done()
		dosenErr = config.DB.Where("id_kelas = ? AND nip = ?", classID, username).First(&dosenPengampu).Error
	}()

	// Goroutine 2: Verify if student is registered in the class
	
	go func() {
		defer wg.Done()
		mahasiswaErr = config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", req.NPM, classID, "aktif").First(&mahasiswaKelas).Error
	}()

	// Goroutine 3: Check if student has already checked in for this meeting
	
	go func() {
		defer wg.Done()
		presensiErr = config.DB.Where("npm = ? AND id_pertemuan = ?", req.NPM, meetingID).First(&existingPresensi).Error
	}()

	// Wait for all validation queries to complete
	wg.Wait()

	// Check lecturer access first
	if dosenErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk kelas ini"})
		return
	}

	// Check student registration
	if mahasiswaErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mahasiswa tidak terdaftar di kelas ini"})
		return
	}

	// Check if student already has attendance (presensiErr == nil means record exists)
	if presensiErr == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Mahasiswa sudah melakukan presensi untuk pertemuan ini"})
		return
	}

	// All validations passed, create attendance record
	presensi := models.Presensi{
		NPM:           req.NPM,
		IDPertemuan:   meetingID,
		WaktuPresensi: time.Now(),
		Status:        "Hadir",
		Catatan:       req.Catatan,
	}

	// Insert attendance record
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