package models

type ProgramStudi struct {
	KodeProdi    string   `gorm:"primaryKey" json:"kode_prodi"`
	NamaProdi    string   `json:"nama_prodi"`
	KodeFakultas string   `json:"kode_fakultas"`
	Fakultas     Fakultas `gorm:"foreignKey:KodeFakultas;references:KodeFakultas"`
}
func (ProgramStudi) TableName() string {
	return "program_studi"
}