package models

type Jadwal struct {
	IDJadwal    int         `gorm:"primaryKey;autoIncrement" json:"id_jadwal"`
	Hari        string      `json:"hari"`
	IDKelas     *int        `json:"id_kelas"`
	IDSesi      *int        `json:"id_sesi"`
	KodeRuangan string      `json:"kode_ruangan"`
	Kelas       Kelas       `json:"kelas,omitempty" gorm:"foreignKey:IDKelas" `
	Sesi        Sesi        `gorm:"foreignKey:IDSesi"`
	Ruangan     Ruangan     `gorm:"foreignKey:KodeRuangan"`
	Pertemuan   []Pertemuan `gorm:"foreignKey:IDJadwal"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}
