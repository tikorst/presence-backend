// package mobile

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/tikorst/presence-backend/config"
// 	"github.com/tikorst/presence-backend/models"
// 	"golang.org/x/crypto/bcrypt"
// )

// // Request struct for login
// type LoginRequest struct {
// 	Username string `json:"username" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// 	DeviceID string `json:"device_id" binding:"required"`
// }

// // Function to handle login requests
// func Login(c *gin.Context) {

// 	// Check if request payload is valid
// 	var req LoginRequest
// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, LoginResponse{Error: true, Message: "Input tidak valid"})
// 		return
// 	}

// 	// Check if username is exists and get user data
// 	var user models.User
// 	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
// 		log.Printf("Login Error: User '%s' not found or DB error: %v", req.Username, err)
// 		c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Username atau password salah"})
// 		return
// 	}

// 	// create a channel for password verification
// 	passVerifyChan := make(chan error, 1)

// 	// go routine to verify password
// 	go func() {
// 		passVerifyChan <- bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
// 	}()

// 	select {

// 		// Check if password verification is not successful
// 		case err := <-passVerifyChan:
// 			if err != nil {
// 				log.Printf("Login Error: Password mismatch for user '%s'", req.Username)
// 				c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Username atau password salah"})
// 				return
// 			}

// 		// Timeout if password verification takes too long
// 		case <-time.After(5 * time.Second):
// 			log.Printf("Login Error: Password verification timeout for user '%s'", req.Username)
// 			c.JSON(http.StatusInternalServerError, LoginResponse{Error: true, Message: "Gagal memverifikasi password (timeout)"})
// 			return
// 	}

// 	// Check if user is a student
// 	if user.TipeUser != "Mahasiswa" {
// 		c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Aplikasi ini hanya untuk Mahasiswa"})
// 		return
// 	}


// 	// Check if the device ID is the same as the one stored in the database
// 	if user.DeviceID != req.DeviceID {

// 		// Check if the device ID was updated within the last 7 days
// 		if user.DeviceIDUpdatedAt != nil && time.Since(*user.DeviceIDUpdatedAt).Hours() < 168 {
// 			c.JSON(http.StatusForbidden, LoginResponse{Error: true, Message: "Pergantian perangkat belum diizinkan"})
// 			return
// 		}

// 		// Try to update the device ID in the database
// 		if err := config.DB.Model(&user).Updates(map[string]interface{}{
// 			"device_id":            req.DeviceID,
// 			"device_id_updated_at": time.Now(),
// 		}).Error; err != nil {

// 			// If the error is due to a duplicate key, return a specific error message
// 			if strings.Contains(err.Error(), "duplicated key") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
// 				c.JSON(http.StatusForbidden, LoginResponse{Error: true, Message: "Perangkat sudah digunakan oleh user lain"})
// 				return
// 			} else {

// 				// Log the error and return a generic error message
// 				log.Printf("Login Error: Failed to update device_id for user '%s': %v", req.Username, err)
// 				c.JSON(http.StatusInternalServerError, LoginResponse{Error: true, Message: "Gagal memperbarui perangkat"})
// 				return
// 			}
// 		}
// 	}

// 	// Generate JWT token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub":       user.Username,
// 		"device_id": req.DeviceID,
// 		"exp":       time.Now().Add(7 * 24 * time.Hour).Unix(),
// 	})

// 	// Sign the token with the secret key
// 	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

// 	// Check if there was an error signing the token
// 	if err != nil {
// 		log.Printf("Login Error: Failed to generate token for user '%s': %v", req.Username, err)
// 		c.JSON(http.StatusInternalServerError, LoginResponse{Error: true, Message: "Gagal membuat token"})
// 		return
// 	}

// 	// Log the user login using go routine
// 	go createUserLog(&user, req.DeviceID, c, true)

// 	// Create the login result
// 	loginResult := &LoginResult{
// 		Name:        user.Nama,
// 		Username:    user.Username,
// 		TokenString: tokenString,
// 	}

// 	// Return the login response
// 	c.JSON(http.StatusOK, LoginResponse{
// 		LoginResult: loginResult,
// 		Error:       false,
// 		Message:     "Login berhasil",
// 	})
// }

// // Function to create a user log entry
// func createUserLog(user *models.User, deviceID string, c *gin.Context, success bool) {

// 	// Get the client IP address and user agent
// 	ip := c.ClientIP()
// 	ua := c.Request.UserAgent()

// 	// parse userid to uint
// 	var userID uint
// 	if user != nil {
// 		userID = uint(user.IDUser)
// 	}

// 	// Create a new user log entry
// 	logEntry := models.UserLog{
// 		IDUser:    userID,
// 		DeviceID:  deviceID,
// 		IPAddress: ip,
// 		UserAgent: ua,
// 		LoginTime: time.Now(),
// 		Success:   success,
// 	}

// 	// Insert the log entry into the database
// 	if err := config.DB.Create(&logEntry).Error; err != nil {
// 		log.Printf("UserLog Error: gagal mencatat log login: %v", err)
// 	}
// }
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

// Request struct for login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}

// Function to handle login requests
func Login(c *gin.Context) {

	// Check if request payload is valid
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{Error: true, Message: "Input tidak valid"})
		return
	}

	// Check if username is exists and get user data
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.Printf("Login Error: User '%s' not found or DB error: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Username atau password salah"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Login Error: Password mismatch for user '%s'", req.Username)
		c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Username atau password salah"})
		return
	}

	// Check if user is a student
	if user.TipeUser != "Mahasiswa" {
		c.JSON(http.StatusUnauthorized, LoginResponse{Error: true, Message: "Aplikasi ini hanya untuk Mahasiswa"})
		return
	}

	// Check if the device ID is the same as the one stored in the database
	if user.DeviceID != req.DeviceID {

		// Check if the device ID was updated within the last 7 days
		if user.DeviceIDUpdatedAt != nil && time.Since(*user.DeviceIDUpdatedAt).Hours() < 168 {
			c.JSON(http.StatusForbidden, LoginResponse{Error: true, Message: "Pergantian perangkat belum diizinkan"})
			return
		}

		// Try to update the device ID in the database
		if err := config.DB.Model(&user).Updates(map[string]interface{}{
			"device_id":            req.DeviceID,
			"device_id_updated_at": time.Now(),
		}).Error; err != nil {

			// If the error is due to a duplicate key, return a specific error message
			if strings.Contains(err.Error(), "duplicated key") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
				c.JSON(http.StatusForbidden, LoginResponse{Error: true, Message: "Perangkat sudah digunakan oleh pengguna lain"})
				return
			} else {

				// Log the error and return a generic error message
				log.Printf("Login Error: Failed to update device_id for user '%s': %v", req.Username, err)
				c.JSON(http.StatusInternalServerError, LoginResponse{Error: true, Message: "Gagal memperbarui perangkat"})
				return
			}
		}
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.Username,
		"device_id": req.DeviceID,
		"exp":       time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	// Check if there was an error signing the token
	if err != nil {
		log.Printf("Login Error: Failed to generate token for user '%s': %v", req.Username, err)
		c.JSON(http.StatusInternalServerError, LoginResponse{Error: true, Message: "Gagal membuat token"})
		return
	}

	// Log the user login using go routine (non-blocking)
	go createUserLog(&user, req.DeviceID, c, true)

	// Create the login result
	loginResult := &LoginResult{
		Name:        user.Nama,
		Username:    user.Username,
		TokenString: tokenString,
	}

	// Return the login response
	c.JSON(http.StatusOK, LoginResponse{
		LoginResult: loginResult,
		Error:       false,
		Message:     "Login berhasil",
	})
}

// Function to create a user log entry
func createUserLog(user *models.User, deviceID string, c *gin.Context, success bool) {

	// Get the client IP address and user agent
	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	// parse userid to uint
	var userID uint
	if user != nil {
		userID = uint(user.IDUser)
	}

	// Create a new user log entry
	logEntry := models.UserLog{
		IDUser:    userID,
		DeviceID:  deviceID,
		IPAddress: ip,
		UserAgent: ua,
		LoginTime: time.Now(),
		Success:   success,
	}

	// Insert the log entry into the database
	if err := config.DB.Create(&logEntry).Error; err != nil {
		log.Printf("UserLog Error: gagal mencatat log login: %v", err)
	}
}