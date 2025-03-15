package web

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/siatma-backend/config"
	"github.com/tikorst/siatma-backend/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

		if user.TipeUser != "Dosen" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Website ini hanya untuk dosen"})
			return
		}
		// Async Password Verification

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		// Cek hasil verifikasi password
		if passErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}

		// Async Update Device ID jika berbeda

		// Generate JWT Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Username,
			"exp": time.Now().Add(1 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token", "err": err})
			return
		}
		c.SetCookie("token", tokenString, 3600, "/", "localhost", true, false)
		// Kirim response ke client
		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"user":    user,
		})
	}
}
