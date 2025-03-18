package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

type QRRequest struct {
	QRCode    string  `json:"qr_code" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	DeviceID  string  `json:"device_id" binding:"required"`
	Timestamp string  `json:"timestamp" binding:"required"`
}

func ValidateQr(c *gin.Context) {
	// Validate the input
	var req QRRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}
	fmt.Println("Location", req.Latitude, req.Longitude)
	// check if the QR code is valid in Redis
	classIDStr, err := config.RedisDB.Get(config.Ctx, req.QRCode).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "QR code tidak valid"})
		return
	}
	// Convert the classID to int
	classID, err := strconv.Atoi(classIDStr)
	// Get the username(npm) from the JWT claims
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonversi classID"})
		return
	}
	claims, _ := c.Get("claims")
	jwtClaims := claims.(jwt.MapClaims)
	username := jwtClaims["sub"].(string)

	var pertemuan models.Pertemuan
	if err := config.DB.Preload("Jadwal.Kelas.MataKuliah").Where("id_pertemuan = ?", classID).First(&pertemuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan informasi kelas", "err": err.Error()})
		return
	}
	// check if user is enrolled in the class
	var mahasiswaKelas models.MahasiswaKelas
	if err := config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", username, pertemuan.Jadwal.Kelas.IDKelas, "aktif").First(&mahasiswaKelas).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu tidak terdaftar di kelas ini"})
		return
	}
	// check if user already attended the class
	if err := config.DB.Where("npm = ? AND id_pertemuan = ?", username, classID).First(&models.Presensi{}).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu sudah melakukan presensi"})
		return
	}
	// Create a new presensi record
	attendance := models.Presensi{
		NPM:           username,
		IDPertemuan:   classID,
		WaktuPresensi: time.Now(),
		Status:        "hadir",
	}
	if err := config.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data presensi", "err": err.Error()})
		return
	}
	// config.DB.Create()
	// fmt.Println("QR code:", qr)
	// config.RedisDB.Get(config.Ctx, qr)
	// // Check if the QR code is valid
	// if qr == "123456" {
	// 	c.JSON(http.StatusOK, gin.H{"message": "QR code is valid"})
	// }
	var message = " Berhasil presensi di kelas " + pertemuan.Jadwal.Kelas.MataKuliah.NamaMatkul + " - " + pertemuan.Jadwal.Kelas.NamaKelas
	c.JSON(http.StatusOK, gin.H{"message": message})
}
