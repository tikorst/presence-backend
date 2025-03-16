package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
	"github.com/tikorst/presence-backend/config"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Ganti sesuai kebutuhan di produksi
	},
}

func GenerateQR(c *gin.Context) {
	classID := c.Param("classID")
	fmt.Println("classID:", classID)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	expired := time.Now().Add(3 * time.Minute)
	for time.Now().Before(expired) {
		qrContent := generateRandomString()
		qr, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
		if err != nil {
			ws.WriteJSON(map[string]string{"error": "failed to generate QR code"})
			continue
		}
		qrBase64 := base64.StdEncoding.EncodeToString(qr)
		config.RedisDB.Set(config.Ctx, qrContent, classID, 17*time.Minute)
		fmt.Println("classID:", classID)
		ws.WriteJSON(map[string]string{"qr": qrBase64})
		time.Sleep(15 * time.Second)
	}
	ws.WriteJSON(map[string]string{"complete": "Stream ended after 3 minutes"})
}

// func generateQR(c *gin.Context) {
// 	classID := c.Param("classID")
// 	qrContent := generateRandomString()

// 	// Simpan QR di Redis dengan TTL 3 menit
// 	redisdb.Client.Set(redisdb.Ctx, classID, qrContent, 3*time.Minute)

// 	// Generate QR Code
// 	qr, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate QR code"})
// 		return
// 	}

// 	c.Data(http.StatusOK, "image/png", qr)
// }

// func validateQR(c *gin.Context) {
// 	var request struct {
// 		QRCode   string `json:"qr_code"`
// 		ClassID  string `json:"class_id"`
// 		Location string `json:"location"`
// 	}
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
// 		return
// 	}

// 	// Cek QR di Redis
// 	storedQR, err := redisdb.Client.Get(redisdb.Ctx, request.ClassID).Result()
// 	if err == redis.Nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "QR code not found"})
// 		return
// 	}

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
// 		return
// 	}

// 	// Validasi QR dan Lokasi
// 	if storedQR != request.QRCode {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid QR code"})
// 		return
// 	}

// 	// Misalnya cek lokasi sesuai dengan parameter
// 	if !isValidLocation(request.Location) {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid location"})
// 		return
// 	}

// 	// Simpan ke database (presensi berhasil)
// 	c.JSON(http.StatusOK, gin.H{"message": "attendance valid"})
// }

// Generate random string untuk QR Code
func generateRandomString() string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, 10)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// // Validasi lokasi (contoh sederhana)
// func isValidLocation(location string) bool {
// 	return location == "valid_location" // Ganti dengan logika validasi lokasi
// }
