package mobile

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
	 "github.com/tikorst/presence-backend/helpers"
)

// QRRequest struct to hold the request data for QR validation
type QRRequest struct {
	QRCode    string  `json:"qr_code" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Timestamp string  `json:"timestamp" binding:"required"`
}

// ValidateQr handles the QR code validation and attendance marking
func ValidateQr(c *gin.Context) {
	// validate request payload
	var req QRRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}
	// Check if qr code exists in Redis
	pertemuanIDStr, err := config.RedisDB.Get(config.Ctx, req.QRCode).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code tidak valid"})
		return
	}
	
	// Convert pertemuanID from string to int
	pertemuanID, err := strconv.Atoi(pertemuanIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonversi classID"})
		return
	}
	// get username and device_id from JWT token
	username, _ := helpers.GetUsername(c)
	device_id, _ := helpers.GetDeviceID(c)

	// Check whether student has already checked in for this meeting
	var existingPresensi models.Presensi
	if err := config.DB.Where("npm = ? AND id_pertemuan = ?", username, pertemuanID).First(&existingPresensi).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu sudah melakukan presensi"})
		return
	}

	// set key for jadwal from pertemuanID
	jadwalKey := fmt.Sprintf("jadwal:%d", pertemuanID)
	var pertemuan models.Pertemuan
	var jadwal models.Jadwal

	// Check redis for jadwal data
	jadwalData, err := config.RedisDB.Get(config.Ctx, jadwalKey).Result()
	if err == nil {
		// If not error, unmarshal jadwal data from Redis
		if err := json.Unmarshal([]byte(jadwalData), &jadwal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pertemuan dari Redis"})
			return
		}
	} else {
		// Get jadwal data from database if not found in Redis
		if err := config.DB.Where("id_pertemuan = ?", pertemuanID).First(&pertemuan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan informasi pertemuan", "err": err.Error()})
			return
		}

		// Get Related Jadwal data
		if err := config.DB.Preload("Ruangan").Preload("Kelas.MataKuliah").
			Where("id_jadwal = ?", pertemuan.IDJadwal).First(&jadwal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan informasi kelas", "err": err.Error()})
			return
		}

		// save jadwal data to Redis with 5 minutes expiration
		if jadwalJSON, err := json.Marshal(jadwal); err == nil {
			config.RedisDB.Set(config.Ctx, jadwalKey, jadwalJSON, 5*time.Minute)
		}
	}

	// Check whether student is registered in the class
	var mahasiswaKelas models.MahasiswaKelas
	if err := config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", username, jadwal.IDKelas, "aktif").First(&mahasiswaKelas).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu tidak terdaftar di kelas ini"})
		return
	}

	// Check if user location is within 30 meters of the class location
	distance := haversine(req.Latitude, req.Longitude, jadwal.Ruangan.Latitude.Float64, jadwal.Ruangan.Longitude.Float64)
	if distance > float64(30) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Lokasi kamu terlalu jauh dari kelas"})
		return
	}

	// create attendance record
	attendance := models.Presensi{
		NPM:           username,
		IDPertemuan:   pertemuanID,
		WaktuPresensi: time.Now(),
		DeviceID:      device_id,
		Status:        "Hadir",
	}

	// Save attendance record to the database
	if err := config.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data presensi"})
		return
	}

	// Return success message
	var message = "Berhasil presensi di kelas " + jadwal.Kelas.MataKuliah.NamaMatkul + " - " + jadwal.Kelas.NamaKelas
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// Haversine function to calculate the distance between two points on the Earth
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // radius of the Earth in meters

	// Convert latitude and longitude from degrees to radians
	dLat := (lat2 - lat1) * math.Pi / 180 
	dLon := (lon2 - lon1) * math.Pi / 180

	
	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
