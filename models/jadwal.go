package models

type Jadwal struct {
	IDJadwal  int         `gorm:"primaryKey;autoIncrement" json:"id_jadwal"`
	Hari      string      `json:"hari"`
	IDKelas   *int        `json:"id_kelas"`
	IDSesi    *int        `json:"id_sesi"`
	IDRuangan *int        `json:"id_ruangan"`
	Kelas     Kelas       `gorm:"foreignKey:IDKelas"`
	Sesi      Sesi        `gorm:"foreignKey:IDSesi"`
	Ruangan   Ruangan     `gorm:"foreignKey:IDRuangan"`
	Pertemuan []Pertemuan `gorm:"foreignKey:IDJadwal" json:"pertemuan"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}
