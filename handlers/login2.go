package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
	"golang.org/x/crypto/bcrypt"
)

func Login2() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Database Ping Test

		// Channel untuk komunikasi antar goroutines
		userChan := make(chan *models.User, 1)
		errChan := make(chan error, 1)

		// Async Query: Fetch User dari Database
		go func() {
			var user models.User
			if err := config.DB.Where("username = ?", "tikorst").First(&user).Error; err != nil {
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
		passChan := make(chan error, 1)
		go func() {
			passErr := bcrypt.CompareHashAndPassword([]byte("$2a$11$PeqzYfKUFiSxNJOvwqlb0OlRc8LSa83suaT9EU9cMAs9fC6LsMRg."), []byte("tikorst"))
			passChan <- passErr
		}()

		// Cek hasil verifikasi password
		if passErr := <-passChan; passErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password incorrect"})
			return
		}

		// Async Check if another user is using the same Device ID
		deviceChan := make(chan bool, 1)
		go func() {
			var otherUser models.User
			err := config.DB.Where("device_id = ? AND id_user != ?", "ABC123", user.IDUser).First(&otherUser).Error
			deviceChan <- err == nil // True jika ada user lain dengan device yang sama
		}()

		// Cek hasil device check
		if <-deviceChan {
			c.JSON(http.StatusForbidden, gin.H{"error": "Device already in use by another user"})
			return
		}

		// Async Update Device ID jika berbeda
		go func() {
			config.DB.Model(user).Where("id_user = ?", user.IDUser).Update("device_id", "ABC1233")
		}()
		// Generate JWT Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.IDUser,
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
