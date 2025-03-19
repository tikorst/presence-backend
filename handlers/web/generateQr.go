package web

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
		return true
	},
}

func GenerateQR(c *gin.Context) {
	// classID := c.Param("classID")
	meetingID := c.Param("meetingID")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	// Channel untuk mendeteksi close
	done := make(chan struct{})

	// Goroutine buat baca pesan supaya server tahu kapan client disconnect
	go func() {
		defer close(done)
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("Client disconnected:", err)
				return
			}
		}
	}()

	expired := time.Now().Add(3 * time.Minute)

	// Goroutine untuk mengirim QR code setiap 15 detik
	go func() {
		for {
			select {
			case <-done:
				log.Println("WebSocket closed, stopping QR generation.")
				return
			default:
				if time.Now().After(expired) {
					ws.WriteJSON(map[string]string{"complete": "Stream ended after 3 minutes"})
					return
				}
				qrContent := generateRandomString()
				qr, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
				if err != nil {
					ws.WriteJSON(map[string]string{"error": "failed to generate QR code"})
					continue
				}

				qrBase64 := base64.StdEncoding.EncodeToString(qr)
				config.RedisDB.Set(config.Ctx, qrContent, meetingID, 17*time.Second)

				// Pastikan WebSocket masih terbuka sebelum mengirim data
				err = ws.WriteJSON(map[string]string{"qr": qrBase64})
				if err != nil {
					log.Println("Failed to send message, closing WebSocket", err)
					return
				}

				time.Sleep(15 * time.Second)
			}
		}
	}()

	// Goroutine untuk mengirim data presensi setiap 5 detik
	// go func() {
	// 	for {
	// 		select {
	// 		case <-done:
	// 			log.Println("WebSocket closed, stopping attendance updates.")
	// 			return
	// 		default:
	// 			var mahasiswaKelas []models.MahasiswaKelas
	// 			if err := config.DB.Preload("Mahasiswa").Preload("Mahasiswa.User").Where("id_kelas = ? AND status = ?", classID, "aktif").Find(&mahasiswaKelas).Order("npm DESC").Error; err != nil {
	// 				c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak ada mahasiswa yang terdaftar di kelas ini"})
	// 				return
	// 			}
	// 			var presensi []models.Presensi
	// 			if err := config.DB.Where("id_pertemuan = ?", meetingID).Find(&presensi).Error; err != nil {
	// 				c.JSON(http.StatusUnauthorized, gin.H{"error": "Belum ada presensi untuk pertemuan ini"})
	// 				return
	// 			}

	// 			attendanceMap := make(map[string]string)
	// 			for _, mhs := range mahasiswaKelas {
	// 				attendanceMap[mhs.Mahasiswa.User.Nama] = "Alpha"
	// 			}

	// 			for _, pres := range presensi {
	// 				attendanceMap[pres.NPM] = "Hadir"
	// 			}

	// 			// Create a response data structure
	// 			responseData := make([]map[string]string, 0)
	// 			for _, mhs := range mahasiswaKelas {
	// 				responseData = append(responseData, map[string]string{
	// 					"nama":   mhs.Mahasiswa.User.Nama,
	// 					"status": attendanceMap[mhs.NPM],
	// 				})
	// 			}

	// 			// Send attendance data
	// 			err = ws.WriteJSON(map[string]interface{}{"attendance": responseData})
	// 			if err != nil {
	// 				log.Println("Failed to send message, closing WebSocket", err)
	// 				return
	// 			}

	// 			time.Sleep(5 * time.Second)
	// 		}
	// 	}
	// }()

	// Block the main goroutine until the WebSocket is closed or expired
	select {
	case <-done:
	case <-time.After(3 * time.Minute):
	}
}

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
