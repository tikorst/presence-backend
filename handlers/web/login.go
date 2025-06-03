package web

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
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Cek apakah user dengan username tersebut ada di database
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Verifikasi apakah user adalah Dosen atau Admin
	if user.TipeUser != "Dosen" && user.TipeUser != "Admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Website ini hanya untuk Dosen dan Admin"})
		return
	}

	// Verifikasi password
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	// Cek hasil verifikasi password
	if passErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}

	// Async Update Device ID jika berbeda

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.Username,
		"role": user.TipeUser,
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	})

	// Sign token dengan secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token", "err": err})
		return
	}
	// Set cookie dengan token JWT
	c.Header("Set-Cookie", "token="+tokenString+"; Path=/; Domain=.tikorst.cloud; Max-Age=3600; Secure; SameSite=None")
	// Kirim response ke client
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"user":    user,
	})
}
