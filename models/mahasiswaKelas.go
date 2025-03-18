package models

type MahasiswaKelas struct {
	NPM       string    `gorm:"primaryKey" json:"npm"`
	IDKelas   int       `gorm:"primaryKey" json:"id_kelas"`
	Status    string    `json:"status"`
	Mahasiswa Mahasiswa `gorm:"foreignKey:NPM;references:NPM"`
}

func (MahasiswaKelas) TableName() string {
	return "mahasiswa_kelas"
}
