package mobile

import (
	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
)

func GetAllGrade(c *gin.Context) {
	username, _ := helpers.GetUsername(c)

	var gradeList []GradeResponse
	err := config.DB.Debug().Table("kelas").
		Joins("JOIN mahasiswa_kelas ON kelas.id_kelas = mahasiswa_kelas.id_kelas").
		Joins("JOIN nilai ON mahasiswa_kelas.id_kelas = nilai.id_kelas AND mahasiswa_kelas.npm = nilai.npm").
		Joins("JOIN mata_kuliah ON kelas.id_matkul = mata_kuliah.id_matkul").
		Where("mahasiswa_kelas.npm = ? AND nilai.nilai_huruf IS NOT NULL", username).
		Select("mahasiswa_kelas.npm, kelas.id_kelas, kelas.id_semester, mata_kuliah.kode_matkul AS kode_matkul, mata_kuliah.nama_matkul AS nama_matkul, mata_kuliah.total_sks AS total_sks, nilai.*").
		Scan(&gradeList).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data nilai"})
		return
	}

	c.JSON(200, gin.H{"error": false, "data": gradeList})
}
