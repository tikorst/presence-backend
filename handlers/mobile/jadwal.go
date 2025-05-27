package mobile

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
	"gorm.io/gorm"
)

type JadwalResponse struct {
	NPM            string             `json:"npm"`
	IDKelas        int                `json:"id_kelas"`
	IDJadwal       int                `json:"id_jadwal"`
	Hari           string             `json:"hari"`
	NamaMataKuliah string             `json:"nama_matkul"`
	NamaKelas      string             `json:"nama_kelas"`
	NamaDosen      string             `json:"nama_dosen"`
	IDRuangan      string             `json:"id_ruangan"`
	Sesi           string             `json:"sesi"`
	Pertemuan      []models.Pertemuan `json:"pertemuan"`
}

func Schedules(c *gin.Context) {
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
	SubQuery := config.DB.
	Table("semester").
	Select("id_semester").
	Order("tahun_ajaran DESC").
	Limit(1)

	if err := config.DB.
		Preload("MataKuliah").
		Preload("DosenPengampu.Dosen").
		Preload("DosenPengampu.Dosen.User").
		Preload("Jadwal", func(db *gorm.DB) *gorm.DB {
			return db.Omit("kelas")
		}).
		Preload("Jadwal.Sesi").
		Preload("Jadwal.Ruangan").
		Preload("Jadwal.Pertemuan").
		Where("id_kelas IN ? AND id_semester = (?)", kelasIDs, SubQuery).Find(&kelas).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data kelas"})
		return
	}
	var jadwalResponse []JadwalResponse
	for _, k := range kelas {
		for _, j := range k.Jadwal {
			jamMasuk := formatTime(j.Sesi.JamMasuk)
			jamKeluar := formatTime(j.Sesi.JamKeluar)
			jadwalResponse = append(jadwalResponse, JadwalResponse{
				NPM:            username,
				IDKelas:        k.IDKelas,
				IDJadwal:       j.IDJadwal,
				Hari:           j.Hari,
				NamaMataKuliah: k.MataKuliah.NamaMatkul,
				NamaKelas:      k.NamaKelas,
				NamaDosen:      k.DosenPengampu[0].Dosen.User.Nama,
				IDRuangan:      j.KodeRuangan,
				Sesi:           jamMasuk + " - " + jamKeluar,
				Pertemuan:      j.Pertemuan,
			})
		}
	}
	fmt.Println("NPM:", username)
	c.JSON(200, gin.H{"erorr": false, "data": jadwalResponse})
}
func formatTime(input string) string {
	parsedTime, err := time.Parse("15:04:05", input)
	if err != nil {
		return input // Return the original input if parsing fails
	}
	return parsedTime.Format("15:04")
}
