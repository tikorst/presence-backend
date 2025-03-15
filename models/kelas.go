package models

type Kelas struct {
	IDKelas    int        `gorm:"primaryKey;autoIncrement" json:"id_kelas"`
	NamaKelas  string     `json:"nama_kelas"`
	KodeMatkul string     `json:"kode_matkul"`
	Kapasitas  *int       `json:"kapasitas"`
	MataKuliah MataKuliah `gorm:"foreignKey:KodeMatkul;references:KodeMatkul"`
}

func (Kelas) TableName() string {
	return "kelas"
}
