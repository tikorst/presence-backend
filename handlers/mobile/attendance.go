package mobile

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

// import (
// 	"fmt"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/tikorst/presence-backend/config"
// 	"github.com/tikorst/presence-backend/models"
// )

// type AttendanceResponse struct {
// 	IDKelas  int `json:"id_kelas"`
// 	IDJadwal int `json:"id_jadwal"`
// }
// type PresensiRingkas struct {
// 	NamaMatkul     string    `json:"nama_matkul"`
// 	NamaKelas      string    `json:"nama_kelas"`
// 	Tanggal        time.Time `json:"tanggal"`
// 	StatusPresensi string    `json:"status"`
// }

// func Attendance(c *gin.Context) {
// 	claims, _ := c.Get("claims")
// 	jwtClaims := claims.(jwt.MapClaims)
// 	username := jwtClaims["sub"].(string)

// 	// Get the semester ID from query parameters
// 	var req AttendanceRequest
// 	c.ShouldBindJSON(&req)
// 	// If semester_id is empty, use the latest semester
// 	if idSemester == 0 {
// 		var latestSemester models.Semester
// 		if err := config.DB.Debug().
// 			Last(&latestSemester).Error; err != nil {
// 			c.JSON(500, gin.H{"error": "Gagal mengambil semester terakhir"})
// 			return
// 		}
// 		fmt.Println("id_semester", latestSemester)
// 		req.IDSemester = latestSemester.IDSemester
// 	}

// 	// Mengambil data kelas yang diambil mahasiswa

// 	var hasil []PresensiRingkas
// 	err := config.DB.Debug().
// 		Table("presensi").
// 		Select(`mata_kuliah.nama_matkul,
// 			kelas.nama_kelas,
// 			pertemuan.tanggal,
// 			presensi.status`).
// 		Joins("JOIN pertemuan ON presensi.id_pertemuan = pertemuan.id_pertemuan").
// 		Joins("JOIN jadwal ON pertemuan.id_jadwal = jadwal.id_jadwal").
// 		Joins("JOIN kelas ON jadwal.id_kelas = kelas.id_kelas").
// 		Joins("JOIN mata_kuliah ON kelas.id_matkul = mata_kuliah.id_matkul").
// 		Where("kelas.id_semester = ? AND presensi.npm = ?", req.IDSemester, username).
// 		Scan(&hasil).Error

// 	if err != nil {
// 		c.JSON(500, gin.H{"error": "Gagal mengambil presensi"})
// 		return
// 	}
// 	fmt.Println("id_semester", req.IDSemester)
// 	c.JSON(200, gin.H{"Error": false, "Message": "Berhasil Mengambil data presensi", "data": hasil})
// }

type KelasResponse struct {
	IDKelas    int           `json:"id_kelas"`
	NamaKelas  string        `json:"nama_kelas"`
	MataKuliah string        `json:"mata_kuliah"`
	IDMatkul   int           `json:"id_matkul"`
	IDSemester int           `json:"id_semester"`
	Presensi   []PresensiRes `json:"presensi"`
}

type PresensiRes struct {
	IDPresensi  int    `json:"id_presensi"`
	IDKelas     int    `json:"id_kelas"`
	IDPertemuan int    `json:"id_pertemuan"`
	PertemuanKe int    `json:"pertemuan_ke"`
	Tanggal     string `json:"tanggal"`
	Status      string `json:"status"`
}

func Attendance(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtClaims := claims.(jwt.MapClaims)
	username := jwtClaims["sub"].(string)
	idSemesterStr := c.Query("id_semester")
	// If semester_id is empty, use the latest semester
	idSemester, err := strconv.Atoi(idSemesterStr)

	if err != nil {
		// kalau gagal parsing, balikin error
		idSemester = 0
	}

	if idSemester == 0 {
		var latestSemester models.Semester
		if err := config.DB.Debug().
			Last(&latestSemester).Error; err != nil {
			c.JSON(500, gin.H{"error": "Gagal mengambil semester terakhir"})
			return
		}
		fmt.Println("id_semester", latestSemester)
		idSemester = latestSemester.IDSemester
	}

	var kelasList []KelasResponse
	config.DB.Table("kelas").
		Joins("JOIN mata_kuliah ON kelas.id_matkul = mata_kuliah.id_matkul").
		Joins("JOIN mahasiswa_kelas ON kelas.id_kelas = mahasiswa_kelas.id_kelas").
		Where("kelas.id_semester = ? AND mahasiswa_kelas.npm = ?", idSemester, username).
		Select("kelas.id_kelas, kelas.nama_kelas, mata_kuliah.nama_matkul as mata_kuliah, kelas.id_matkul, kelas.id_semester").
		Scan(&kelasList)

	var presensiList []PresensiRes
	config.DB.Table("presensi").Debug().
		Joins("JOIN pertemuan ON presensi.id_pertemuan = pertemuan.id_pertemuan").
		Joins("JOIN jadwal ON pertemuan.id_jadwal = jadwal.id_jadwal").
		Joins("JOIN kelas ON jadwal.id_kelas = kelas.id_kelas").
		Where("kelas.id_semester = ? AND presensi.npm = ?", idSemester, username).
		Select("presensi.id_presensi, kelas.id_kelas, pertemuan.id_pertemuan, pertemuan.pertemuan_ke, pertemuan.tanggal, presensi.status").
		Order("pertemuan.tanggal").
		Scan(&presensiList)

	// Kelompokkan presensi berdasarkan id_kelas
	presensiMap := make(map[int][]PresensiRes)
	for _, p := range presensiList {
		presensiMap[p.IDKelas] = append(presensiMap[p.IDKelas], p)
	}

	// Inject ke dalam masing-masing kelas
	for i := range kelasList {
		kelasList[i].Presensi = presensiMap[kelasList[i].IDKelas]
	}

	c.JSON(200, gin.H{"erorr": false, "data": kelasList, "npm": username})
}
