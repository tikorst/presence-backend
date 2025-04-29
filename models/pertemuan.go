package models

import "time"

type Pertemuan struct {
	IDPertemuan int       `gorm:"primaryKey;autoIncrement" json:"id_pertemuan"`
	Tanggal     time.Time `json:"tanggal"`
	Status      string    `json:"status"`
	PertemuanKe string    `json:"pertemuan_ke"`
	IDJadwal    *int      `json:"id_jadwal"`
	Jadwal      Jadwal    `gorm:"foreignKey:IDJadwal;references:IDJadwal"`
}

func (Pertemuan) TableName() string {
	return "pertemuan"
}
