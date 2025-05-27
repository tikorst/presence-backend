package mobile

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

type GradeResponse struct {
	NPM        string  `json:"npm"`
	IDKelas    int     `json:"id_kelas"`
	IDSemester int     `json:"id_semester"`
	NamaKelas  string  `json:"nama_kelas"`
	NamaMatkul string  `json:"nama_matkul"`
	TotalSKS   int     `json:"total_sks"`
	IDNilai    int     `json:"id_nilai"`
	NilaiUTS   float64 `json:"nilai_uts"`
	NilaiUAS   float64 `json:"nilai_uas"`
	NilaiHuruf string  `json:"nilai_huruf"`
	Bobot      float64 `json:"bobot"`
}

func Grade(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtClaims := claims.(jwt.MapClaims)
	username := jwtClaims["sub"].(string)
	idSemesterStr := c.Query("id_semester")

	idSemester, err := strconv.Atoi(idSemesterStr)

	if err != nil {
		// kalau gagal parsing, balikin error
		idSemester = 0
	}

	if idSemester == 0 {
		var latestSemester models.Semester
		if err := config.DB.
			Last(&latestSemester).Error; err != nil {
			c.JSON(500, gin.H{"error": "Gagal mengambil semester terakhir"})
			return
		}
		fmt.Println("id_semester", latestSemester)
		idSemester = latestSemester.IDSemester
	}

	var gradeList []GradeResponse
	err = config.DB.Table("kelas").
		Joins("JOIN mahasiswa_kelas ON kelas.id_kelas = mahasiswa_kelas.id_kelas").
		Joins("JOIN nilai ON mahasiswa_kelas.id_kelas = nilai.id_kelas AND mahasiswa_kelas.npm = nilai.npm").
		Joins("JOIN mata_kuliah ON kelas.id_matkul = mata_kuliah.id_matkul").
		Where("kelas.id_semester = ? AND mahasiswa_kelas.npm = ?", idSemester, username).
		Select("mahasiswa_kelas.npm, kelas.id_kelas, kelas.id_semester, kelas.nama_kelas, mata_kuliah.nama_matkul AS nama_matkul, mata_kuliah.total_sks AS total_sks, nilai.*").
		Scan(&gradeList).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data nilai"})
		return
	}

	c.JSON(200, gin.H{"error": false, "data": gradeList})
}
