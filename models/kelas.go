package models

type Kelas struct {
	IDKelas       int             `gorm:"primaryKey;autoIncrement" json:"id_kelas"`
	NamaKelas     string          `json:"nama_kelas"`
	IDMatkul      string          `json:"id_matkul"`
	IDSemester    int             `json:"id_semester"`
	Kapasitas     *int            `json:"kapasitas"`
	MataKuliah    MataKuliah      `gorm:"foreignKey:IDMatkul;references:IDMatkul"`
	DosenPengampu []DosenPengampu `gorm:"foreignKey:IDKelas;references:IDKelas"`
	Jadwal        []Jadwal        `gorm:"foreignKey:IDKelas;references:IDKelas"`
	Semester      Semester        `gorm:"foreignKey:IDSemester;references:IDSemester"`
}

func (Kelas) TableName() string {
	return "kelas"
}
