package models

type Kelas struct {
	IDKelas       int             `gorm:"primaryKey;autoIncrement" json:"id_kelas"`
	NamaKelas     string          `json:"nama_kelas"`
	IDMatkul      string          `json:"id_matkul"`
	Kapasitas     *int            `json:"kapasitas"`
	MataKuliah    MataKuliah      `gorm:"foreignKey:IDMatkul;references:IDMatkul"`
	DosenPengampu []DosenPengampu `gorm:"foreignKey:IDKelas;references:IDKelas"`
	Jadwal        []Jadwal        `gorm:"foreignKey:IDKelas;references:IDKelas"`
}

func (Kelas) TableName() string {
	return "kelas"
}
