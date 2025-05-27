package mobile

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"


	"cloud.google.com/go/storage"
)

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

func Profile(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtClaims := claims.(jwt.MapClaims)
	username := jwtClaims["sub"].(string)

	mhs := models.Mahasiswa{}
	if err := config.DB.Debug().
		Model(&mhs).Where("npm = ?", username).
		Preload("User").
		Preload("ProgramStudi").
		First(&mhs).Error; err != nil {
		c.JSON(500, gin.H{"error": true, "message": "Gagal mengambil data mahasiswa"})
		return
	}


	type serviceAccount struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}
	var sa serviceAccount
	if err := json.Unmarshal([]byte(config.JSONCreds), &sa); err != nil {
		c.JSON(500, gin.H{"error": true, "message": "GCP_SERVICE_ACCOUNT_JSON env var is empty"})
		return
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: sa.ClientEmail,
		PrivateKey:     []byte(sa.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(1 * time.Hour),
	}

	bucketName := os.Getenv("BUCKET_NAME")
	objectName := "profile-picture/" + mhs.NPM + ".jpg"

	url, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		c.JSON(500, gin.H{"error": true, "message": "failed to generate signed URL: %v"})
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
		URL:          url,
	}
	c.JSON(200, gin.H{"error": false, "data": res})
}
