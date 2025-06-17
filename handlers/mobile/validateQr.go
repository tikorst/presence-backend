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

type QRRequest struct {
	QRCode    string  `json:"qr_code" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Timestamp string  `json:"timestamp" binding:"required"`
}

func ValidateQr(c *gin.Context) {
	// Validasi input
	var req QRRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}
	// Validasi QR Code
	pertemuanIDStr, err := config.RedisDB.Get(config.Ctx, req.QRCode).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "QR code tidak valid"})
		return
	}
	// Validasi pertemuanID
	pertemuanID, err := strconv.Atoi(pertemuanIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonversi classID"})
		return
	}
	// Ambil claims dari context
	// Pastikan JWT token valid
	username, _ := helpers.GetUsername(c)
	device_id, _ := helpers.GetDeviceID(c)

	// Check apakah mahasiswa sudah melakukan presensi
	var existingPresensi models.Presensi
	if err := config.DB.Where("npm = ? AND id_pertemuan = ?", username, pertemuanID).First(&existingPresensi).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu sudah melakukan presensi"})
		return
	}

	// Check apakah mahasiswa terdaftar di kelas tersebut

	// Ambil data pertemuan dari database
	jadwalKey := fmt.Sprintf("jadwal:%d", pertemuanID)
	var pertemuan models.Pertemuan
	var jadwal models.Jadwal

	jadwalData, err := config.RedisDB.Get(config.Ctx, jadwalKey).Result()
	if err == nil {
		// Data dari Redis
		if err := json.Unmarshal([]byte(jadwalData), &jadwal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pertemuan dari Redis"})
			return
		}
	} else {
		// Ambil dari DB
		if err := config.DB.Where("id_pertemuan = ?", pertemuanID).First(&pertemuan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan informasi pertemuan", "err": err.Error()})
			return
		}
		if err := config.DB.Preload("Ruangan").Preload("Kelas.MataKuliah").
			Where("id_jadwal = ?", pertemuan.IDJadwal).First(&jadwal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan informasi kelas", "err": err.Error()})
			return
		}

		// Simpan ke Redis
		if jadwalJSON, err := json.Marshal(jadwal); err == nil {
			config.RedisDB.Set(config.Ctx, jadwalKey, jadwalJSON, 5*time.Minute)
		}
	}

	var mahasiswaKelas models.MahasiswaKelas
	if err := config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", username, jadwal.IDKelas, "aktif").First(&mahasiswaKelas).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu tidak terdaftar di kelas ini"})
		return
	}

	// Check apakah lokasi presensi sesuai dengan lokasi pertemuan

	distance := haversine(req.Latitude, req.Longitude, jadwal.Ruangan.Latitude.Float64, jadwal.Ruangan.Longitude.Float64)
	if distance > float64(30) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Lokasi kamu terlalu jauh dari kelas"})
		return
	}

	attendance := models.Presensi{
		NPM:           username,
		IDPertemuan:   pertemuanID,
		WaktuPresensi: time.Now(),
		DeviceID:      device_id,
		Status:        "Hadir",
	}
	if err := config.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data presensi"})
		return
	}
	var message = "Berhasil presensi di kelas " + jadwal.Kelas.MataKuliah.NamaMatkul + " - " + jadwal.Kelas.NamaKelas
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Radius bumi dalam meter
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
