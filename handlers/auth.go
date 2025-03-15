package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"github.com/tikorst/siatma-backend/config"
	"github.com/tikorst/siatma-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Async Password Verification

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		// Cek hasil verifikasi password
		if passErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}

		// Async Update Device ID jika berbeda
		if user.DeviceID != req.DeviceID {
			deviceChan := make(chan error, 1)
			go func() {
				// var otherUser models.User
				err := config.DB.Model(user).Update("device_id", req.DeviceID).Where("id_user = ?", user.IDUser).Error
				deviceChan <- err // True jika ada user lain dengan device yang sama
			}()

			// Cek hasil device check

			if deviceErr := <-deviceChan; deviceErr != nil {
				if gorm.ErrDuplicatedKey == deviceErr {
					c.JSON(http.StatusForbidden, gin.H{"error": "Device already in use by another user"})
					return
				} else {
					c.JSON(http.StatusForbidden, gin.H{"error": "Unknown Error"})
					return
				}
			}
		}

		// Generate JWT Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Username,
			"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token", "err": err})
			return
		}

		// Kirim response ke client
		c.JSON(http.StatusOK, gin.H{
			"token_string": tokenString,
		})
	}
}
