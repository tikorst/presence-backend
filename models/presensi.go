package models

import "time"

type Presensi struct {
	PresensiID    int       `gorm:"primaryKey;autoIncrement" json:"presensi_id"`
	NPM           string    `json:"npm"`
	IDPertemuan   int       `json:"id_pertemuan"`
	WaktuPresensi time.Time `json:"waktu_presensi"`
	Status        string    `json:"status"`
	Catatan       string    `json:"catatan" gorm:"type:text"`
	Mahasiswa     Mahasiswa `gorm:"foreignKey:NPM"`
	Pertemuan     Pertemuan `gorm:"foreignKey:IDPertemuan"`
}

func (Presensi) TableName() string {
	return "presensi"
}
