package models

// AsistenDosen represents the asisten_dosen table
type AsistenDosen struct {
	NPM     string `gorm:"primaryKey" json:"npm"`
	IDKelas int    `gorm:"primaryKey" json:"id_kelas"`
	Status  string `json:"status"`
}

// Dosen represents the dosen table

// DosenPengampu represents the dosen_pengampu table

// Fakultas represents the fakultas table

// Jadwal represents the jadwal table

// Kelas represents the kelas table

// Mahasiswa represents the mahasiswa table

// MahasiswaKelas represents the mahasiswa_kelas table

// MataKuliah represents the mata_kuliah table

// Pertemuan represents the pertemuan table

// Presensi represents the presensi table

// ProgramStudi represents the program_studi table

// RefreshToken represents the refresh_token table

// Ruangan represents the ruangan table

// Sesi represents the sesi table

// User represents the users table
