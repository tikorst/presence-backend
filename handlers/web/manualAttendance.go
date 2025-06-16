package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

type ManualAttendanceRequest struct {
	NPM     string `json:"npm" binding:"required"`
	Catatan string `json:"catatan"`
}

func ManualAttendance(c *gin.Context) {
	classIDStr := c.Param("classID")
	meetingIDStr := c.Param("meetingID")
	claims, _ := c.Get("claims")
	jwtClaims := claims.(jwt.MapClaims)
	username := jwtClaims["sub"].(string)
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
		return
	}
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}
	// Get authenticated user from cookies/session

	// Verifikasi apakah user adalah dosen yang mengampu kelas ini
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

	// Verifikasi apakah mahasiswa terdaftar di kelas ini
	var mahasiswaKelas models.MahasiswaKelas
	if err := config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", req.NPM, classID, "aktif").First(&mahasiswaKelas).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mahasiswa tidak terdaftar di kelas ini"})
		return
	}

	// Cek apakah sudah ada presensi untuk pertemuan ini
	var existingPresensi models.Presensi
	if err := config.DB.Where("npm = ? AND id_pertemuan = ?", req.NPM, meetingID).First(&existingPresensi).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Mahasiswa sudah melakukan presensi untuk pertemuan ini"})
		return
	}

	// Insert data presensi
	presensi := models.Presensi{
		NPM:           req.NPM,
		IDPertemuan:   meetingID,
		WaktuPresensi: time.Now(),
		Status:        "Hadir",
		Catatan:       req.Catatan,
	}

	if err := config.DB.Create(&presensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data presensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Presensi berhasil ditambahkan",
		"data":    presensi,
	})
}
