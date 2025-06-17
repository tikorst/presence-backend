package models

import "time"

type Presensi struct {
	IDPresensi    int       `gorm:"primaryKey;autoIncrement" json:"id_presensi"`
	NPM           string    `json:"npm"`
	IDPertemuan   int       `json:"id_pertemuan"`
	WaktuPresensi time.Time `json:"waktu_presensi"`
	Status        string    `json:"status"`
	DeviceID      string    `json:"device_id"`
	Catatan       string    `json:"catatan" gorm:"type:text"`
	Mahasiswa     Mahasiswa `gorm:"foreignKey:NPM"`
	Pertemuan     Pertemuan `gorm:"foreignKey:IDPertemuan"`
}

func (Presensi) TableName() string {
	return "presensi"
}
