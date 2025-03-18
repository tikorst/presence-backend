package models

import "time"

type Pertemuan struct {
	IDPertemuan int       `gorm:"primaryKey;autoIncrement" json:"id_pertemuan"`
	Tanggal     time.Time `json:"tanggal"`
	Status      string    `json:"status"`
	KodeQR      string    `json:"kode_qr"`
	IDJadwal    *int      `json:"id_jadwal"`
	Jadwal 		Jadwal    `gorm:"foreignKey:IDJadwal;references:IDJadwal" json:"jadwal"`
}

func (Pertemuan) TableName() string {
	return "pertemuan"
}
