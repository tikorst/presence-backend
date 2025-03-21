package models

type Fakultas struct {
	KodeFakultas string `gorm:"primaryKey" json:"kode_fakultas"`
	NamaFakultas string `json:"nama_fakultas"`
}

func (Fakultas) TableName() string {
	return "fakultas"
}
