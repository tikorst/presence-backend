package mobile

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/helpers"
	"github.com/tikorst/presence-backend/models"
)

type KelasResponse struct {
	IDKelas    int           `json:"id_kelas"`
	NamaKelas  string        `json:"nama_kelas"`
	MataKuliah string        `json:"mata_kuliah"`
	IDMatkul   int           `json:"id_matkul"`
	IDSemester int           `json:"id_semester"`
	Presensi   []PresensiRes `json:"presensi" gorm:"-"`
}

type PresensiRes struct {
	IDPresensi  int    `json:"id_presensi"`
	IDKelas     int    `json:"id_kelas"`
	IDPertemuan int    `json:"id_pertemuan"`
	PertemuanKe int    `json:"pertemuan_ke"`
	Tanggal     string `json:"tanggal"`
	Status      string `json:"status"`
}

func GetAttendance(c *gin.Context) {
	username, _ := helpers.GetUsername(c)
	idSemesterStr := c.Query("id_semester")
	// If semester_id is empty, use the latest semester
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

	var kelasList []KelasResponse
	config.DB.Table("kelas").
		Joins("JOIN mata_kuliah ON kelas.id_matkul = mata_kuliah.id_matkul").
		Joins("JOIN mahasiswa_kelas ON kelas.id_kelas = mahasiswa_kelas.id_kelas").
		Where("kelas.id_semester = ? AND mahasiswa_kelas.npm = ?", idSemester, username).
		Select("kelas.id_kelas, kelas.nama_kelas, mata_kuliah.nama_matkul as mata_kuliah, kelas.id_matkul, kelas.id_semester").
		Scan(&kelasList)

	var presensiList []PresensiRes
	today := time.Now().Format("2006-01-02")

	config.DB.Table("pertemuan").
		Select(`
            COALESCE(presensi.id_presensi, 0) as id_presensi,
            kelas.id_kelas,
            pertemuan.id_pertemuan,
            pertemuan.pertemuan_ke,
            pertemuan.tanggal,
            COALESCE(presensi.status, 'Alpha') as status
        `).
		Joins("JOIN jadwal ON pertemuan.id_jadwal = jadwal.id_jadwal").
		Joins("JOIN kelas ON jadwal.id_kelas = kelas.id_kelas").
		Joins("LEFT JOIN presensi ON presensi.id_pertemuan = pertemuan.id_pertemuan AND presensi.npm = ?", username).
		Where("kelas.id_semester = ? AND (pertemuan.tanggal < ? OR pertemuan.status = ?)", idSemester, today, "selesai").
		Order("pertemuan.tanggal").
		Scan(&presensiList)

	// Kelompokkan presensi berdasarkan id_kelas
	for i := range presensiList {
		if presensiList[i].Status != "" {
			presensiList[i].Status = strings.ToUpper(presensiList[i].Status[:1]) + strings.ToLower(presensiList[i].Status[1:])
		}
	}
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
