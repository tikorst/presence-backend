package models

type ProgramStudi struct {
	IDProdi    int      `gorm:"primaryKey;autoIncrement" json:"id_prodi"`
	KodeProdi  string   `json:"kode_prodi"`
	NamaProdi  string   `json:"nama_prodi"`
	IDFakultas *int     `json:"id_fakultas"`
	Fakultas   Fakultas `gorm:"foreignKey:IDFakultas"`
}
