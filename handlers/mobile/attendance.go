package mobile

import (
	"strconv"
	"strings"
	"sync"
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
	// Get the username from the context
	username, _ := helpers.GetUsername(c)

	// Get the semester_id from the query parameters
	idSemesterStr := c.Query("id_semester")
	idSemester, err := strconv.Atoi(idSemesterStr)

	if err != nil {
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

	// Set today's date in the format "YYYY-MM-DD"
	today := time.Now().Format("2006-01-02")

	// Channel untuk menampung hasil dari goroutine
	type QueryResult struct {
		KelasList    []KelasResponse
		PresensiList []PresensiRes
		Error        error
	}

	// Channel untuk komunikasi antar goroutine
	var presensiList []PresensiRes
	var kelasList []KelasResponse
	var presensiErr, kelasErr error
	var wg sync.WaitGroup

	// Goroutine untuk query kelas
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		
		kelasErr = config.DB.Table("kelas").
			Joins("JOIN mata_kuliah ON kelas.id_matkul = mata_kuliah.id_matkul").
			Joins("JOIN mahasiswa_kelas ON kelas.id_kelas = mahasiswa_kelas.id_kelas").
			Where("kelas.id_semester = ? AND mahasiswa_kelas.npm = ?", idSemester, username).
			Select("kelas.id_kelas, kelas.nama_kelas, mata_kuliah.nama_matkul as mata_kuliah, kelas.id_matkul, kelas.id_semester").
			Scan(&kelasList).Error
	}()

	// Goroutine untuk query presensi
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		
		presensiErr = config.DB.Table("pertemuan").
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
			Where("kelas.id_semester = ? AND (pertemuan.tanggal < ? OR pertemuan.status = ? OR presensi.id_presensi IS NOT NULL)", idSemester, today, "selesai").
			Order("pertemuan.tanggal").
			Scan(&presensiList).Error
	}()

	// Tunggu semua goroutine selesai
	wg.Wait()

	if( kelasErr != nil || presensiErr != nil) {
		c.JSON(500, gin.H{"error": "Gagal mengambil data kelas atau presensi"})
		return
	}
	// Process presensi status formatting in parallel
	// Kelompokkan presensi berdasarkan id_kelas
	for i := range presensiList {
		if presensiList[i].Status != "" {
			presensiList[i].Status = strings.ToUpper(presensiList[i].Status[:1]) + strings.ToLower(presensiList[i].Status[1:])
		}
	}

	// Buat map untuk mengelompokkan presensi berdasarkan ID kelas
	presensiMap := make(map[int][]PresensiRes)
	for _, p := range presensiList {
		presensiMap[p.IDKelas] = append(presensiMap[p.IDKelas], p)
	}

	// Inject ke dalam masing-masing kelas
	for i := range kelasList {
		kelasList[i].Presensi = presensiMap[kelasList[i].IDKelas]
	}

	// Return the list of classes with attendance records
	c.JSON(200, gin.H{"erorr": false, "data": kelasList, "npm": username})
}