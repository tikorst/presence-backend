package web

import (
	"net/http"

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
	if err := config.DB.
		Last(&latestSemester).Error; err != nil {
		// If there's an error in fetching the latest semester, return an error response
		c.JSON(500, gin.H{"error": "Gagal mengambil semester terakhir"})
		return
	}

	// Query to get the list of classes for the user in the latest semester
	// This query joins the dosen_pengampu table to filter classes by the user's NIP
	if err := config.DB.
		Joins("JOIN dosen_pengampu ON dosen_pengampu.id_kelas = kelas.id_kelas").
		Preload("MataKuliah").
		Preload("DosenPengampu").
		Where("dosen_pengampu.nip = ? AND id_semester = ?", username, latestSemester.IDSemester).
		Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch classes"})
		return
	}

	// Query to get the user details based on the username
	if err := config.DB.
		Where("username = ?", username).
		First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// If everything is successful, return the classes and user details
	c.JSON(http.StatusOK, gin.H{"classes": classes, "user": user})
}
