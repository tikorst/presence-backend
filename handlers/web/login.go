package web

import (
	"encoding/base64"
	"math/rand"
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
	tokenString, csrfToken, err := createJWTWithCSRF(user.Username, user.TipeUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	// Set cookie dengan token JWT
	c.Header("Set-Cookie", "token="+tokenString+"; Path=/; Domain=.tikorst.cloud; Max-Age=3600; HttpOnly; SameSite=Lax; Secure")

	// Kirim response ke client
	c.JSON(http.StatusOK, gin.H{
		"message":    "Login berhasil",
		"csrf_token": csrfToken,
		"user":       user,
	})
}
func generateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Create JWT dengan CSRF token embedded
func createJWTWithCSRF(username, role string) (string, string, error) {
	csrfToken, err := generateCSRFToken()
	if err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        username,
		"role":       role,
		"csrf_token": csrfToken, // Embed CSRF token
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", "", err
	}

	return tokenString, csrfToken, nil
}
