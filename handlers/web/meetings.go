package web

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

// Struct untuk response jadwal
type JadwalResponse struct {
	IDJadwal    int                `json:"id_jadwal"`
	Hari        string             `json:"hari"`
	IDKelas     int                `json:"id_kelas"`
	IDSesi      int                `json:"id_sesi"`
	KodeRuangan string             `json:"kode_ruangan"`
	Sesi        models.Sesi        `json:"Sesi"`
	Ruangan     models.Ruangan     `json:"Ruangan"`
	Pertemuan   []models.Pertemuan `json:"Pertemuan"`
}

func GetMeetings(c *gin.Context) {
	// Ambil classID dari URL
	classID := c.Param("classID")

	var jadwal []models.Jadwal
	var pertemuan []models.Pertemuan
	var jadwalErr, pertemuanErr error
	var wg sync.WaitGroup

	// Parallel queries menggunakan WaitGroup
	wg.Add(2)

	// Goroutine 1: Ambil data jadwal berdasarkan kelas
	go func() {
		defer wg.Done()
		jadwalErr = config.DB.
			Joins("Ruangan").
			Joins("Sesi").
			Where("jadwal.id_kelas = ?", classID).
			Find(&jadwal).Error
	}()

	// Goroutine 2: Ambil semua pertemuan untuk kelas ini
	go func() {
		defer wg.Done()
		pertemuanErr = config.DB.
			Where("id_jadwal IN (?)",
				config.DB.
					Table("jadwal").
					Select("id_jadwal").
					Where("id_kelas = ?", classID),
			).
			Find(&pertemuan).Error
	}()

	// Wait for both queries to complete
	wg.Wait()

	// Check for errors
	if jadwalErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}

	if pertemuanErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meetings"})
		return
	}

	// Jika tidak ada jadwal ditemukan
	if len(jadwal) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No schedules found for this class"})
		return
	}

	// Buat map untuk mengelompokkan pertemuan berdasarkan id_jadwal
	pertemuanMap := make(map[int][]models.Pertemuan)
	for _, p := range pertemuan {
		pertemuanMap[p.IDJadwal] = append(pertemuanMap[p.IDJadwal], p)
	}

	// Transformasi hasil ke struct response
	var jadwalResponse []JadwalResponse
	for _, j := range jadwal {
		// Ambil pertemuan yang sesuai dengan jadwal ini
		jadwalPertemuan := pertemuanMap[j.IDJadwal]
		if jadwalPertemuan == nil {
			jadwalPertemuan = []models.Pertemuan{} // Empty slice jika tidak ada pertemuan
		}

		jadwalResponse = append(jadwalResponse, JadwalResponse{
			IDJadwal:    j.IDJadwal,
			Hari:        j.Hari,
			IDKelas:     j.IDKelas,
			IDSesi:      j.IDSesi,
			KodeRuangan: j.KodeRuangan,
			Sesi:        j.Sesi,
			Ruangan:     j.Ruangan,
			Pertemuan:   jadwalPertemuan,
		})
	}

	// Kirim response ke client
	c.JSON(http.StatusOK, gin.H{
		"status": "Berhasil",
		"data":   jadwalResponse,
		"count":  len(jadwalResponse),
	})
}