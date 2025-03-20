package mobile

import (
	"encoding/json"
	"fmt"
	"math"
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
	fmt.Println("Location data:", req.Latitude, req.Longitude)
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
	fmt.Println("Username:", username)
	pertemuanChan := make(chan *models.Pertemuan)
	mahasiswaKelasChan := make(chan *models.MahasiswaKelas)
	presensiChan := make(chan error)
	distanceChan := make(chan float64)

	// Check apakah mahasiswa sudah melakukan presensi
	go func() {
		if err := config.DB.Where("npm = ? AND id_pertemuan = ?", username, classID).First(&models.Presensi{}).Error; err == nil {
			presensiChan <- fmt.Errorf("Kamu sudah melakukan presensi")
			return
		}
		presensiChan <- nil
	}()

	if err := <-presensiChan; err != nil {
		fmt.Println("Presensi error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check apakah mahasiswa terdaftar di kelas tersebut
	go func() {
		var mahasiswaKelas models.MahasiswaKelas
		if err := config.DB.Where("npm = ? AND id_kelas = ? AND status = ?", username, classID, "aktif").First(&mahasiswaKelas).Error; err != nil {
			mahasiswaKelasChan <- nil
			return
		}
		mahasiswaKelasChan <- &mahasiswaKelas
	}()

	mahasiswaKelas := <-mahasiswaKelasChan
	if mahasiswaKelas == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kamu tidak terdaftar di kelas ini"})
		return
	}

	// Ambil data pertemuan dari database
	pertemuanKey := fmt.Sprintf("pertemuan:%d", classID)
	go func() {
		pertemuanData, err := config.RedisDB.Get(config.Ctx, pertemuanKey).Result()
		var pertemuan models.Pertemuan
		if err == nil {
			// Data found in Redis, unmarshal it
			fmt.Println("data ketemu", pertemuanData)
			if err := json.Unmarshal([]byte(pertemuanData), &pertemuan); err != nil {
				fmt.Println("error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pertemuan dari Redis"})

				return
			}
		} else {
			// Data not found in Redis, query the database
			fmt.Println("data tidak ketemu")
			if err := config.DB.Preload("Jadwal.Ruangan").Preload("Jadwal.Kelas.MataKuliah").Where("id_pertemuan = ?", classID).First(&pertemuan).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan informasi kelas", "err": err.Error()})
				return
			}
			// Store the data in Redis with a TTL of 5 minutes
			pertemuanJSON, err := json.Marshal(pertemuan)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pertemuan untuk Redis"})
				return
			}
			config.RedisDB.Set(config.Ctx, pertemuanKey, pertemuanJSON, 5*time.Minute)
		}
		pertemuanChan <- &pertemuan
	}()

	pertemuan := <-pertemuanChan
	if pertemuan == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan informasi kelas"})
		return
	}

	// Check apakah lokasi presensi sesuai dengan lokasi pertemuan
	fmt.Println("Pertemuan location:", pertemuan.Jadwal.Ruangan)
	go func() {
		distance := haversine(req.Latitude, req.Longitude, pertemuan.Jadwal.Ruangan.Latitude.Float64, pertemuan.Jadwal.Ruangan.Longitude.Float64)
		distanceChan <- distance
	}()

	distance := <-distanceChan
	fmt.Println("Distance:", float64(distance))
	if distance > float64(30) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Lokasi kamu terlalu jauh dari lokasi kelas"})
		return
	}

	// config.DB.Create()
	// fmt.Println("QR code:", qr)
	// config.RedisDB.Get(config.Ctx, qr)
	// // Check if the QR code is valid
	// if qr == "123456" {
	// 	c.JSON(http.StatusOK, gin.H{"message": "QR code is valid"})
	// }

	attendance := models.Presensi{
		NPM:           username,
		IDPertemuan:   classID,
		WaktuPresensi: time.Now(),
		Status:        "hadir",
	}
	if err := config.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data presensi"})
		return
	}
	var message = " Berhasil presensi di kelas " + pertemuan.Jadwal.Kelas.MataKuliah.NamaMatkul + " - " + pertemuan.Jadwal.Kelas.NamaKelas
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
