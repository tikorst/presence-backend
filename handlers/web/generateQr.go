package web

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"crypto/rand"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
	"gorm.io/gorm"
)

// websocket upgrader to upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Function to handle WebSocket connection and generate QR codes
func GenerateQR(c *gin.Context) {

	// get meetingID from URL parameters
	meetingID := c.Param("meetingID")
	classID := c.Param("classID")

	username, _ := helpers.GetUsername(c)

	// convert classID and meetingID from string to int
	classIDInt, err := strconv.Atoi(classID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
		return
	}

	// convert meetingID from string to int
	meetingIDInt, err := strconv.Atoi(meetingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	var (
		wg            sync.WaitGroup
		authErr       error
		matchErr      error
		pertemuan     models.Pertemuan
		dosenPengampu models.DosenPengampu
	)

	wg.Add(2)

	// Goroutine 1: check if user is a lecturer for the class
	go func() {
		defer wg.Done()
		authErr = config.DB.
			Where("id_kelas = ? AND nip = ?", classIDInt, username).
			First(&dosenPengampu).Error
	}()

	// Goroutine 2: check if the meeting matches the class
	go func() {
		defer wg.Done()
		matchErr = config.DB.
			Joins("JOIN jadwal ON jadwal.id_jadwal = pertemuan.id_jadwal").
			Where("pertemuan.id_pertemuan = ? AND jadwal.id_kelas = ?", meetingIDInt, classIDInt).
			First(&pertemuan).Error
	}()

	wg.Wait()

	// validate the results of the goroutines
	if authErr != nil {
		if authErr == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengampu kelas ini"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi akses kelas"})
		}
		return
	}

	if matchErr != nil {
		if matchErr == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pertemuan tidak ditemukan untuk kelas ini"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi data pertemuan"})
		}
		return
	}

	// uptrade HTTP connection to WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// make sure to close the WebSocket connection when done
	defer ws.Close()

	// make a channel to signal when the WebSocket is closed
	done := make(chan struct{})

	// Goroutine to handle if the WebSocket connection is closed
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

	// Set expiration time for QR code generation
	// QR codes will be generated for 3 minutes
	expired := time.Now().Add(3 * time.Minute)

	// Goroutine to generate QR codes every 15 seconds
	// Loops to continuously generate QR codes until the WebSocket is closed or the time expires
	go func() {
		for {
			select {
			case <-done: // If the done channel is closed, exit the goroutine
				log.Println("WebSocket closed, stopping QR generation.")
				return

			default: // default case to generate QR codes every 15 seconds
				if time.Now().After(expired) {
					ws.WriteJSON(map[string]string{"complete": "Stream ended after 3 minutes"})
					return
				}

				// Generate a random string for the QR code content then encode it to a QR code
				qrContent, _ := generateRandomString()
				qr, err := qrcode.Encode(qrContent, qrcode.Medium, 256)

				// If there is an error generating the QR code, send an error message to the WebSocket
				if err != nil {
					ws.WriteJSON(map[string]string{"error": "failed to generate QR code"})
					continue
				}

				// Store it in Redis AND check for errors, qrContent is the key and meetingID is the value
				redisValue := fmt.Sprintf("%s:%s", meetingID, classID)
				err = config.RedisDB.Set(config.Ctx, qrContent, redisValue, 17*time.Second).Err()
				if err != nil {
					// Log the error for debugging and try again in the next loop
					log.Println("Error setting key in Redis:", err)
					// Optionally, inform the client that there was a temporary issue
					ws.WriteJSON(map[string]string{"error": "internal server error, retrying"})
					continue
				}

				qrBase64 := base64.StdEncoding.EncodeToString(qr) // Convert the QR code to a base64 string
				// Send the QR code as a JSON message to the WebSocket
				err = ws.WriteJSON(map[string]string{"qr": qrBase64})

				// If there is an error sending the message, log the error and close the WebSocket
				if err != nil {
					log.Println("Failed to send message, closing WebSocket", err)
					return
				}

				// Sleep for 15 seconds before generating the next QR code
				time.Sleep(15 * time.Second)
			}
		}
	}()

	// Select to wait for the done channel or a timeout
	select {
	case <-done:
	case <-time.After(3 * time.Minute):
	}
}

// Generate a random string of 10 characters for QR code content
func generateRandomString() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err // Handle error if reading from the random reader fails
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
