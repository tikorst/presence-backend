package models

// AsistenDosen represents the asisten_dosen table
type AsistenDosen struct {
	NPM     string `gorm:"primaryKey" json:"npm"`
	IDKelas int    `gorm:"primaryKey" json:"id_kelas"`
	Status  string `json:"status"`
}

