package models

type Fakultas struct {
	IDFakultas   int    `gorm:"primaryKey" json:"id_fakultas"`
	KodeFakultas string `gorm:"primaryKey" json:"kode_fakultas"`
	NamaFakultas string `json:"nama_fakultas"`
}
