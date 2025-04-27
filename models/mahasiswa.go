package models

type Mahasiswa struct {
	NPM                     string       `gorm:"primaryKey;column:npm" json:"npm"`
	IDUser                  int          `json:"id_user"`
	KodeProdi               string       `json:"kode_prodi"`
	TahunMasuk              string       `json:"tahun_masuk"`
	Status                  string       `json:"status"`
	DosenPembimbingAkademik string       `json:"dosen_pembimbing_akademik"`
	User                    User         `gorm:"foreignKey:NPM;references:Username"`
	ProgramStudi            ProgramStudi `gorm:"foreignKey:KodeProdi;references:KodeProdi"`
	Dosen                   Dosen        `gorm:"foreignKey:DosenPembimbingAkademik;references:NIP"`
}

func (Mahasiswa) TableName() string {
	return "mahasiswa"
}
