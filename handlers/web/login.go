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

// LoginRequest struct to bind the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {

	// Validate the request payload
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the user exists in the database
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Check if the user type is either "Dosen" or "Admin"
	if user.TipeUser != "Dosen" && user.TipeUser != "Admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Website ini hanya untuk Dosen dan Admin"})
		return
	}

	// Verify the password
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	
	// Check if password verification failed
	if passErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}

	// Generate JWT Token
	tokenString, csrfToken, err := createJWTWithCSRF(user.Username, user.TipeUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	// Set cookie with the JWT token
	c.Header("Set-Cookie", "token="+tokenString+"; Path=/; Domain=.tikorst.cloud; Max-Age=3600; HttpOnly; SameSite=Lax; Secure")

	// Set CSRF token in the response header
	c.JSON(http.StatusOK, gin.H{
		"message":    "Login berhasil",
		"csrf_token": csrfToken,
		"user":       user,
	})
}

// function to generate a CSRF token
func generateCSRFToken() (string, error) {

	// Generate a random CSRF token using crypto/rand
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode the bytes to a URL-safe base64 string
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Create JWT dengan CSRF token embedded
func createJWTWithCSRF(username, role string) (string, string, error) {

	// Generate CSRF token
	csrfToken, err := generateCSRFToken()
	if err != nil {
		return "", "", err
	}

	// Create JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        username,
		"role":       role,
		"csrf_token": csrfToken, // Embed CSRF token
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
	})

	// Sign the token with the secret key from environment variable
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", "", err
	}
	

	// Return the signed token string and CSRF token
	return tokenString, csrfToken, nil
}
