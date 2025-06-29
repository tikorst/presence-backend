package web

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
)

func GetClasses(c *gin.Context) {

	// Get the username from the context
	username, _ := helpers.GetUsername(c)

	// initialize variables
	var classes []models.Kelas
	user := models.User{}
	var latestSemester models.Semester

	// Query to get the latest semester

	// Check Redis cache first
	cachedSemester, err := config.RedisDB.Get(config.Ctx, "latest_semester").Result()

	if err != nil {
		// Not found in cache - query database
		if err := config.DB.
			Last(&latestSemester).Error; err != nil {
			c.JSON(500, gin.H{"error": "Gagal mengambil semester terakhir"})
			return
		}
		semesterJSON, _ := json.Marshal(latestSemester)
		// Store in Redis for next time
		config.RedisDB.Set(config.Ctx, "latest_semester", semesterJSON, 24*time.Hour)
	} else {
		// Found in cache
		if err := json.Unmarshal([]byte(cachedSemester), &latestSemester); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pertemuan dari Redis"})
			return
		}
	}

	// Use goroutines for concurrent database queries
	var wg sync.WaitGroup
	var classesErr, userErr error

	// Goroutine for fetching classes
	wg.Add(2)
	go func() {
		defer wg.Done()
		classesErr = config.DB.
			Joins("JOIN dosen_pengampu ON dosen_pengampu.id_kelas = kelas.id_kelas").
			Joins("MataKuliah").
			Where("dosen_pengampu.nip = ? AND id_semester = ?", username, latestSemester.IDSemester).
			Find(&classes).Error
	}()

	// Goroutine for fetching user details
	go func() {
		defer wg.Done()
		userErr = config.DB.
			Where("username = ?", username).
			First(&user).Error
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Check for errors
	if classesErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch classes"})
		return
	}

	if userErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// If everything is successful, return the classes and user details
	c.JSON(http.StatusOK, gin.H{"classes": classes, "user": user})
}