package models

import "database/sql"

type Ruangan struct {
	KodeRuangan string          `gorm:"primaryKey;unique" json:"kode_ruangan"`
	Status      string          `json:"status"`
	Kapasitas   *int            `json:"kapasitas"`
	Latitude    sql.NullFloat64 `gorm:"type:decimal(10,8)"`
	Longitude   sql.NullFloat64 `gorm:"type:decimal(11,8)"`
}

func (Ruangan) TableName() string {
	return "ruangan"
}
