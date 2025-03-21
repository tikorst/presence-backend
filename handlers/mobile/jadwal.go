package mobile

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
	"gorm.io/gorm"
)

func Jadwal(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtClaims := claims.(jwt.MapClaims)
	username := jwtClaims["sub"].(string)

	// Mengambil data kelas yang diambil mahasiswa
	var mahasiswaKelas []models.MahasiswaKelas
	if err := config.DB.Where("npm = ? AND status = 'aktif'", username).Find(&mahasiswaKelas).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data kelas"})
		return
	}

	var kelasIDs []int
	for _, k := range mahasiswaKelas {
		kelasIDs = append(kelasIDs, k.IDKelas)
	}

	// Mengambil data jadwal berdasarkan kelas yang diambil mahasiswa
	var jadwal []models.Jadwal
	if err := config.DB.
		Preload("Kelas.MataKuliah").
		Preload("Kelas.MataKuliah").
		Preload("Sesi").
		Preload("Ruangan").
		Preload("Pertemuan").Where("id_kelas IN ?", kelasIDs).Find(&jadwal).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data jadwal"})
		return
	}

	var kelas []models.Kelas
	if err := config.DB.
		Preload("MataKuliah").
		Preload("DosenPengampu.Dosen").
		Preload("Jadwal", func(db *gorm.DB) *gorm.DB {
			return db.Omit("kelas")
		}).
		Preload("Jadwal.Sesi").
		Preload("Jadwal.Ruangan").
		Preload("Jadwal.Pertemuan").
		Where("id_kelas IN ?", kelasIDs).Find(&kelas).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data kelas"})
		return
	}

	c.JSON(200, gin.H{"erorr": false, "data": kelas})
}
