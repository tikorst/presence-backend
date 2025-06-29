package mobile

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
)

type GradeResponse struct {
	NPM        string   `json:"npm"`
	IDKelas    int      `json:"id_kelas"`
	IDSemester int      `json:"id_semester"`
	KodeMatkul string   `json:"kode_matkul"`
	NamaMatkul string   `json:"nama_matkul"`
	TotalSKS   int      `json:"total_sks"`
	IDNilai    int      `json:"id_nilai"`
	NilaiUTS   *float64 `json:"nilai_uts"`
	NilaiUAS   *float64 `json:"nilai_uas"`
	NilaiHuruf *string  `json:"nilai_huruf"`
	Bobot      *float64 `json:"bobot"`
}

func GetGrade(c *gin.Context) {

	// Get the username from the context
	username, _ := helpers.GetUsername(c)

	// Get the semester_id from the query parameters
	idSemesterStr := c.Query("id_semester")

	// Convert the semester_id to an integer
	idSemester, err := strconv.Atoi(idSemesterStr)

	// If there's an error in conversion, set idSemester to 0
	if err != nil {
		// kalau gagal parsing, balikin error
		idSemester = 0
	}

	// If idSemester is 0, fetch the latest semester
	if idSemester == 0 {
		// Check Redis cache first
		cachedID, err := config.RedisDB.Get(config.Ctx, "latest_semester_id").Result()

		if err != nil {
			// Not found in cache - query database
			var latestSemester models.Semester
			if err := config.DB.
				Last(&latestSemester).Error; err != nil {
				c.JSON(500, gin.H{"error": "Gagal mengambil semester terakhir"})
				return
			}

			idSemester = latestSemester.IDSemester

			// Store in Redis for next time
			config.RedisDB.Set(config.Ctx, "latest_semester_id", idSemester, 24*time.Hour)
		} else {
			// Found in cache
			idSemester, _ = strconv.Atoi(cachedID)
		}
	}

	// Query to get the list of grades for the user in the specified semester
	var gradeList []GradeResponse
	err = config.DB.Table("kelas").
		Joins("JOIN mahasiswa_kelas ON kelas.id_kelas = mahasiswa_kelas.id_kelas").
		Joins("JOIN nilai ON mahasiswa_kelas.id_kelas = nilai.id_kelas AND mahasiswa_kelas.npm = nilai.npm").
		Joins("JOIN mata_kuliah ON kelas.id_matkul = mata_kuliah.id_matkul").
		Where("kelas.id_semester = ? AND mahasiswa_kelas.npm = ?", idSemester, username).
		Select("mahasiswa_kelas.npm, kelas.id_kelas, kelas.id_semester, mata_kuliah.kode_matkul AS kode_matkul, mata_kuliah.nama_matkul AS nama_matkul, mata_kuliah.total_sks AS total_sks, nilai.*").
		Scan(&gradeList).Error

	// Check for errors in the query
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data nilai"})
		return
	}

	// If no grades found, return an empty list
	c.JSON(200, gin.H{"error": false, "data": gradeList})
}
