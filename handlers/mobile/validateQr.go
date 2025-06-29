package mobile

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
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
	redisValueStr, err := config.RedisDB.Get(config.Ctx, req.QRCode).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code tidak valid"})
		return
	}
	parts := strings.Split(redisValueStr, ":")
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format QR code tidak valid"})
		return
	}

	meetingID := parts[0]
	classID := parts[1]
	// Convert pertemuanID from string to int
	pertemuanID, err := strconv.Atoi(meetingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonversi pertemuanID"})
		return
	}
	// Convert classID from string to int
	kelasID, err := strconv.Atoi(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonversi classID"})
		return
	}
	// get username and device_id from JWT token
	username, _ := helpers.GetUsername(c)
	device_id, _ := helpers.GetDeviceID(c)

	var (
		kelas               models.MahasiswaKelas
		jadwal              models.Jadwal
		jadwalErr, kelasErr error
	)

	var wg sync.WaitGroup
	wg.Add(2)
	// Check whether student is registered in the class
	go func() {
		defer wg.Done()
		kelasErr = config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", username, kelasID, "aktif").First(&kelas).Error

	}()

	// set key for jadwal from pertemuanID
	go func() {
		defer wg.Done()
		jadwalKey := fmt.Sprintf("jadwal:%d", pertemuanID)
		jadwalData, err := config.RedisDB.Get(config.Ctx, jadwalKey).Result()
		if err == nil {
			jadwalErr = json.Unmarshal([]byte(jadwalData), &jadwal)
			return
		}

		// Fallback ke DB
		jadwalErr = config.DB.
			Joins("JOIN pertemuan ON pertemuan.id_jadwal = jadwal.id_jadwal").
			Preload("Ruangan").
			Preload("Kelas.MataKuliah").
			Where("pertemuan.id_pertemuan = ?", pertemuanID).
			First(&jadwal).Error
		if jadwalErr == nil {
			if jadwalJSON, err := json.Marshal(jadwal); err == nil {
				config.RedisDB.Set(config.Ctx, jadwalKey, jadwalJSON, 5*time.Minute)
			}
		}
	}()
	wg.Wait()

	if kelasErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu tidak terdaftar di kelas ini"})
		return
	}
	if jadwalErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan informasi kelas"})
		return
	}
	// Check whether student is registered in the class

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
		if strings.Contains(err.Error(), "duplicated key") || strings.Contains(err.Error(), "UNIQUE") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu sudah melakukan presensi"})
		return
	}
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
