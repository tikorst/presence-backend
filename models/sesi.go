package models

type Sesi struct {
	IDSesi    int    `gorm:"primaryKey;autoIncrement" json:"id_sesi"`
	NoSesi    *int   `json:"no_sesi"`
	JamMasuk  string `json:"jam_masuk" gorm:"type:time"`
	JamKeluar string `json:"jam_keluar" gorm:"type:time"`
}

func (Sesi) TableName() string {
	return "sesi"
}
