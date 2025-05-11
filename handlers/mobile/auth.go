package mobile

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Input tidak valid"})
			return
		}

		// Channel untuk komunikasi antar goroutines
		userChan := make(chan *models.User, 1)
		errChan := make(chan error, 1)

		// Async Query: Fetch User dari Database
		go func() {
			var user models.User
			if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
				errChan <- err
				return
			}
			userChan <- &user
		}()

		// Wait for result dari user query
		var user *models.User
		select {
		case user = <-userChan:
			// User ditemukan, lanjut proses
		case <-errChan:
			c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Username atau password salah"})
			return
		}

		// Async Password Verification

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		// Cek hasil verifikasi password
		if passErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Username atau password salah"})
		}

		// Cek device_id
		if user.DeviceID != req.DeviceID {
			// Cek apakah device_id sudah dipakai user lain
			var otherUserCount int64
			if err := config.DB.Model(&models.User{}).Where("device_id = ? AND username != ?", req.DeviceID, req.Username).Count(&otherUserCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Gagal memeriksa perangkat"})
				return
			}
			if otherUserCount > 0 {
				c.JSON(http.StatusForbidden, gin.H{"error": true, "message": "Perangkat sudah digunakan oleh user lain"})
				return
			}

			// Cek waktu terakhir ganti device
			if user.DeviceIDUpdatedAt != nil && time.Since(*user.DeviceIDUpdatedAt).Hours() < 24 {
				c.JSON(http.StatusForbidden, gin.H{"error": true, "message": "Perangkat baru hanya bisa digunakan setelah 24 jam. Lakukan reset device."})
				return
			}

			// Update device_id dan last_device_change
			if err := config.DB.Model(user).Updates(map[string]interface{}{
				"device_id":            req.DeviceID,
				"device_id_updated_at": time.Now(),
			}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Gagal memperbarui perangkat"})
				return
			}
		}

		// Generate JWT Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Username,
			"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Gagal membuat token"})
			return
		}
		loginResult := &LoginResult{
			Name:        user.Nama,
			Username:    user.Username,
			TokenString: tokenString,
		}

		loginResponse := &LoginResponse{
			LoginResult: loginResult,
			Error:       false,
			Message:     "Login berhasil",
		}
		// Kirim response ke client
		c.JSON(http.StatusOK,
			loginResponse,
		)
	}
}
