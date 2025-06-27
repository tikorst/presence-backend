package mobile

import (
	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
)

func GetAllGrade(c *gin.Context) {

	// Get the username from the context
	username, _ := helpers.GetUsername(c)
	
	var gradeList []GradeResponse

	// Query to get all grades for the user
	err := config.DB.Table("kelas").
		Joins("JOIN mahasiswa_kelas ON kelas.id_kelas = mahasiswa_kelas.id_kelas").
		Joins("JOIN nilai ON mahasiswa_kelas.id_kelas = nilai.id_kelas AND mahasiswa_kelas.npm = nilai.npm").
		Joins("JOIN mata_kuliah ON kelas.id_matkul = mata_kuliah.id_matkul").
		Where("mahasiswa_kelas.npm = ? AND nilai.nilai_huruf IS NOT NULL", username).
		Select("mahasiswa_kelas.npm, kelas.id_kelas, kelas.id_semester, mata_kuliah.kode_matkul AS kode_matkul, mata_kuliah.nama_matkul AS nama_matkul, mata_kuliah.total_sks AS total_sks, nilai.*").
		Scan(&gradeList).Error

	// Check for errors in the query
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data nilai"})
		return
	}

	// Return the list of grades
	c.JSON(200, gin.H{"error": false, "data": gradeList})
}
