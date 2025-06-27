package mobile

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
	

	"cloud.google.com/go/storage"
)

// Custom response struct for profile data
type ProfileResponse struct {
	NPM          string `json:"npm"`
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	NamaProdi    string `json:"nama_prodi"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Telepon      string `json:"telepon"`
	URL          string `json:"url"`
}

// Service account struct for GCP credentials
type serviceAccount struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

// Function to get the profile of a user
func GetProfile(c *gin.Context) {

	// Get the username from the context
	username, _ := helpers.GetUsername(c)

	// Query to get the Mahasiswa data based on the username
	mhs := models.Mahasiswa{}
	if err := config.DB.
		Model(&mhs).Where("npm = ?", username).
		Preload("User").
		Preload("ProgramStudi").
		First(&mhs).Error; err != nil {
		c.JSON(500, gin.H{"error": true, "message": "Gagal mengambil data mahasiswa"})
		return
	}

	// Check if the GCP service account JSON properly set
	var sa serviceAccount
	if err := json.Unmarshal([]byte(config.JSONCreds), &sa); err != nil {
		c.JSON(500, gin.H{"error": true, "message": "GCP_SERVICE_ACCOUNT_JSON is not set or invalid"})
		return
	}

	// Create a signed URL for the profile picture
	opts := &storage.SignedURLOptions{
		GoogleAccessID: sa.ClientEmail,
		PrivateKey:     []byte(sa.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(1 * time.Hour),
	}

	// get the bucket name from environment variable
	bucketName := os.Getenv("BUCKET_NAME")

	// Construct the object name for the profile picture
	objectName := "profile-picture/" + mhs.NPM + ".jpg"

	// Generate the signed URL and handle any errors
	url, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		c.JSON(500, gin.H{"error": true, "message": "failed to generate signed URL: %v"})
		return
	}

	// Create the response struct with the profile data
	res := ProfileResponse{
		NPM:          mhs.NPM,
		Nama:         mhs.User.Nama,
		Email:        mhs.User.Email,
		NamaProdi:    mhs.ProgramStudi.NamaProdi,
		TempatLahir:  mhs.User.TempatLahir,
		TanggalLahir: mhs.User.TanggalLahir.Format("02-01-2006"),
		Telepon:      mhs.User.NoTelepon,
		URL:          url,
	}

	// Return the profile data as JSON
	c.JSON(200, gin.H{"error": false, "data": res})
}
