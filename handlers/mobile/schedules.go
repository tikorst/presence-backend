package mobile

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
)

// ScheduleResponse struct to format the response for schedules
type ScheduleResponse struct {
	NPM        string `json:"npm"`
	IDKelas    int    `json:"id_kelas"`
	IDJadwal   int    `json:"id_jadwal"`
	Hari       string `json:"hari"`
	NamaMatkul string `json:"nama_matkul"`
	NamaKelas  string `json:"nama_kelas"`
	NamaDosen  string `json:"nama_dosen"`
	IDRuangan  string `json:"id_ruangan"`
	Sesi       string `json:"sesi"`
}

// GetSchedule retrieves the schedule for a student
func GetSchedules(c *gin.Context) {
	username, exists := helpers.GetUsername(c)
	if exists != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "User not found in token"})
		return
	}

	idSemester, err := getLatestSemesterID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Gagal mengambil semester terakhir"})
		return
	}

	var jadwalData []ScheduleResponse
	var jadwalErr error
	// var pertemuan []models.Pertemuan

	// Fetch schedule data for the student
	jadwalErr = config.DB.
		Table("mahasiswa_kelas").
		Select(`
			mahasiswa_kelas.npm,
			kelas.id_kelas,
			jadwal.id_jadwal,
			jadwal.hari,
			mata_kuliah.nama_matkul,
			kelas.nama_kelas,
			users.nama as nama_dosen,
			jadwal.kode_ruangan as id_ruangan,
			CONCAT(sesi.jam_masuk, ' - ', sesi.jam_keluar) as sesi
		`).
		Joins("JOIN kelas ON kelas.id_kelas = mahasiswa_kelas.id_kelas").
		Joins("JOIN mata_kuliah ON mata_kuliah.id_matkul = kelas.id_matkul").
		Joins("JOIN jadwal ON jadwal.id_kelas = kelas.id_kelas").
		Joins("JOIN sesi ON sesi.id_sesi = jadwal.id_sesi").
		Joins("JOIN ruangan ON ruangan.kode_ruangan = jadwal.kode_ruangan").
		Joins("LEFT JOIN dosen_pengampu ON dosen_pengampu.id_kelas = kelas.id_kelas").
		Joins("LEFT JOIN dosen ON dosen.nip = dosen_pengampu.nip").
		Joins("LEFT JOIN users ON users.username = dosen.nip").
		Where("mahasiswa_kelas.npm = ? AND mahasiswa_kelas.status = 'aktif' AND kelas.id_semester = ?", username, idSemester).
		Scan(&jadwalData).Error

	// Error handling
	if jadwalErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to fetch schedule data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "data": jadwalData})
}

// getLatestSemesterID retrieves the latest semester ID with Redis caching
func getLatestSemesterID() (int, error) {
	// Check Redis cache first
	cachedID, err := config.RedisDB.Get(config.Ctx, "latest_semester_id").Result()
	if err != nil {
		// Not found in cache - query database
		var latestSemester models.Semester
		if err := config.DB.Last(&latestSemester).Error; err != nil {
			return 0, err
		}

		idSemester := latestSemester.IDSemester

		// Store in Redis for next time
		config.RedisDB.Set(config.Ctx, "latest_semester_id", idSemester, 24*time.Hour)

		return idSemester, nil
	}

	// Found in cache
	idSemester, err := strconv.Atoi(cachedID)
	if err != nil {
		return 0, err
	}

	return idSemester, nil
}
