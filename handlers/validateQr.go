package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

type QRRequest struct {
	QRCode   string `json:"qr_code"`
	Location string `json:"location"`
	DeviceID string `json:"device_id"`
}

func ValidateQr(c *gin.Context) {
	// Validate the input
	var req QRRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// check if the QR code is valid in Redis
	classIDStr, err := config.RedisDB.Get(config.Ctx, req.QRCode).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid QR code"})
		return
	}
	// Convert the classID to int
	classID, err := strconv.Atoi(classIDStr)
	// Get the username(npm) from the JWT claims
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert classID to int"})
		return
	}
	claims, _ := c.Get("claims")
	jwtClaims := claims.(jwt.MapClaims)
	username := jwtClaims["sub"].(string)

	// check if user is enrolled in the class
	var mahasiswaKelas models.MahasiswaKelas
	if err := config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", username, classID, "aktif").First(&mahasiswaKelas).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not enrolled in the class or status is not active"})
		return
	}
	// check if user already attended the class
	if err := config.DB.Where("npm = ? AND id_pertemuan = ?", username, classID).First(&models.Presensi{}).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User already attended the class"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record attendance", "err": err.Error()})
		return
	}
	// config.DB.Create()
	// fmt.Println("QR code:", qr)
	// config.RedisDB.Get(config.Ctx, qr)
	// // Check if the QR code is valid
	// if qr == "123456" {
	// 	c.JSON(http.StatusOK, gin.H{"message": "QR code is valid"})
	// }
	c.JSON(http.StatusOK, gin.H{"message": "Presensi berhasil"})
}
