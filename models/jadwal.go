package models

type Jadwal struct {
	IDJadwal    int         `gorm:"primaryKey;autoIncrement" json:"id_jadwal"`
	Hari        string      `json:"hari"`
	IDKelas     int        `json:"id_kelas"`
	IDSesi      int        `json:"id_sesi"`
	KodeRuangan string      `json:"kode_ruangan"`
	Kelas       Kelas       `json:"kelas" gorm:"foreignKey:IDKelas;references:IDKelas"`
	Sesi        Sesi        `gorm:"foreignKey:IDSesi;references:IDSesi"`
	Ruangan     Ruangan     `gorm:"foreignKey:KodeRuangan;references:KodeRuangan"`
	Pertemuan   []Pertemuan `gorm:"foreignKey:IDJadwal;references:IDJadwal"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}
