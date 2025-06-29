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

// Your existing structs remain the same
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

type serviceAccount struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

type urlResult struct {
	URL string
	Err error
}

// Original function with goroutine (concurrent)
func GetProfile(c *gin.Context) {

	username, _ := helpers.GetUsername(c)
	urlChan := make(chan urlResult)

	var sa serviceAccount
	if err := json.Unmarshal([]byte(config.JSONCreds), &sa); err != nil {
		c.JSON(500, gin.H{"error": true, "message": "GCP_SERVICE_ACCOUNT_JSON is not set or invalid"})
		return
	}

	// Start goroutine for signed URL generation
	go func() {
		opts := &storage.SignedURLOptions{
			GoogleAccessID: sa.ClientEmail,
			PrivateKey:     []byte(sa.PrivateKey),
			Method:         "GET",
			Expires:        time.Now().Add(1 * time.Hour),
		}

		bucketName := os.Getenv("BUCKET_NAME")
		objectName := "profile-picture/" + username + ".jpg"

		url, err := storage.SignedURL(bucketName, objectName, opts)

		urlChan <- urlResult{URL: url, Err: err}
	}()

	// Fetch database data concurrently
	mhs := models.Mahasiswa{}
	if err := config.DB.
		Model(&mhs).Where("npm = ?", username).
		Joins("User").
		Joins("ProgramStudi").
		First(&mhs).Error; err != nil {
		c.JSON(500, gin.H{"error": true, "message": "Gagal mengambil data mahasiswa"})
		return
	}

	// Wait for signed URL result
	result := <-urlChan

	if result.Err != nil {
		c.JSON(500, gin.H{"error": true, "message": "failed to generate signed URL"})
		return
	}

	res := ProfileResponse{
		NPM:          mhs.NPM,
		Nama:         mhs.User.Nama,
		Email:        mhs.User.Email,
		NamaProdi:    mhs.ProgramStudi.NamaProdi,
		TempatLahir:  mhs.User.TempatLahir,
		TanggalLahir: mhs.User.TanggalLahir.Format("02-01-2006"),
		Telepon:      mhs.User.NoTelepon,
		URL:          result.URL,
	}

	c.JSON(200, gin.H{"error": false, "data": res})
}
