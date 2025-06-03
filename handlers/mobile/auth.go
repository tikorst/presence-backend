package mobile

import (
	"log"
	"net/http"
	"os"
	"strings"
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

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{Error: true, Message: "Input tidak valid"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.Printf("Login Error: User '%s' not found or DB error: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Username atau password salah"})
		return
	}

	passVerifyChan := make(chan error, 1)

	go func() {
		passVerifyChan <- bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	}()

	select {
	case err := <-passVerifyChan:
		if err != nil {
			log.Printf("Login Error: Password mismatch for user '%s'", req.Username)
			c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Username atau password salah"})
			return
		}
	case <-time.After(5 * time.Second):
		log.Printf("Login Error: Password verification timeout for user '%s'", req.Username)
		c.JSON(http.StatusInternalServerError, LoginResponse{Error: true, Message: "Gagal memverifikasi password (timeout)"})
		return
	}

	if user.TipeUser != "Mahasiswa" {
		c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Aplikasi ini hanya untuk Mahasiswa"})
		return
	}

	if user.DeviceID != req.DeviceID {
		if user.DeviceIDUpdatedAt != nil && time.Since(*user.DeviceIDUpdatedAt).Hours() < 24 {
			c.JSON(http.StatusForbidden, LoginResponse{Error: true, Message: "Perangkat baru hanya bisa digunakan setelah 24 jam. Lakukan reset device."})
			return
		}

		if err := config.DB.Model(&user).Updates(map[string]interface{}{
			"device_id":            req.DeviceID,
			"device_id_updated_at": time.Now(),
		}).Error; err != nil {
			if strings.Contains(err.Error(), "duplicated key") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
				c.JSON(http.StatusForbidden, LoginResponse{Error: true, Message: "Perangkat sudah digunakan oleh user lain"})
				return
			} else {
				log.Printf("Login Error: Failed to update device_id for user '%s': %v", req.Username, err)
				c.JSON(http.StatusInternalServerError, LoginResponse{Error: true, Message: "Gagal memperbarui perangkat"})
				return
			}
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Printf("Login Error: Failed to generate token for user '%s': %v", req.Username, err)
		c.JSON(http.StatusInternalServerError, LoginResponse{Error: true, Message: "Gagal membuat token"})
		return
	}

	loginResult := &LoginResult{
		Name:        user.Nama,
		Username:    user.Username,
		TokenString: tokenString,
	}

	c.JSON(http.StatusOK, LoginResponse{
		LoginResult: loginResult,
		Error:       false,
		Message:     "Login berhasil",
	})
}
