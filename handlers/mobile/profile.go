package mobile

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

type ProfileResponse struct {
	NPM          string `json:"npm"`
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	NamaProdi    string `json:"nama_prodi"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Telepon      string `json:"telepon"`
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

	res := ProfileResponse{
		NPM:          mhs.NPM,
		Nama:         mhs.User.Nama,
		Email:        mhs.User.Email,
		NamaProdi:    mhs.ProgramStudi.NamaProdi,
		TempatLahir:  mhs.User.TempatLahir,
		TanggalLahir: mhs.User.TanggalLahir.Format("02-01-2006"),
		Telepon:      mhs.User.NoTelepon,
	}
	c.JSON(200, gin.H{"error": false, "data": res})
}
