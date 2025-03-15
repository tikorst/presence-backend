package models

type Ruangan struct {
	IDRuangan   int    `gorm:"primaryKey;autoIncrement" json:"id_ruangan"`
	KodeRuangan string `json:"kode_ruangan"`
	Status      string `json:"status"`
	Kapasitas   *int   `json:"kapasitas"`
}

func (Ruangan) TableName() string {
	return "ruangan"
}
