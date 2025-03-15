package models

type Mahasiswa struct {
	NPM                     string       `gorm:"primaryKey" json:"npm"`
	IDUser                  *int         `json:"id_user"`
	IDProdi                 *int         `json:"id_prodi"`
	TahunMasuk              string       `json:"tahun_masuk"`
	Status                  string       `json:"status"`
	DosenPembimbingAkademik string       `json:"dosen_pembimbing_akademik"`
	User                    User         `gorm:"foreignKey:NPM;references:Username"`
	UserReference           User         `gorm:"foreignKey:IDUser"`
	ProgramStudi            ProgramStudi `gorm:"foreignKey:IDProdi"`
	Dosen                   Dosen        `gorm:"foreignKey:DosenPembimbingAkademik;references:NIP"`
}
