package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

type PresensiRequest struct {
	NPM      string `json:"npm" binding:"required"`
	QRCode   string `json:"qr_code" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}

func ScanQR() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req PresensiRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Cek QR code valid
		var pertemuan models.Pertemuan
		err := config.DB.First(&pertemuan, "kode_qr = ? AND tanggal = ?", req.QRCode, time.Now().Format("2006-01-02")).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or expired QR code"})
			return
		}
		// check device_id
		if err := config.DB.First(&models.User{}, "npm = ? AND device_id = ?", req.NPM, req.DeviceID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Device not authorized"})
			return
		}

		presensi := models.Presensi{
			NPM:           req.NPM,
			IDPertemuan:   pertemuan.IDPertemuan,
			WaktuPresensi: time.Now(),
			Status:        "hadir",
		}

		if err := config.DB.Create(&presensi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record attendance", "err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Presensi berhasil"})
	}
}
