package models

type Dosen struct {
	NIP       string `gorm:"primaryKey" grom:"unique" json:"nip"`
	IDUser    int    `json:"id_user"`
	IDProdi   int    `json:"id_prodi"`
	KodeDosen string `json:"kode_dosen"`
	Status    string `json:"status"`
	// User          User         `gorm:"foreignKey:NIP;references:Username"`
	// UserReference User         `gorm:"foreignKey:IDUser"`
	// ProgramStudi  ProgramStudi `gorm:"foreignKey:IDProdi"`
}

func (Dosen) TableName() string {
	return "dosen"
}
