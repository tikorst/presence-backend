package models

type DosenPengampu struct {
	NIP     string `gorm:"primaryKey;column:nip" json:"nip"`
	IDKelas int    `gorm:"primaryKey" json:"id_kelas"`
	Status  string `json:"status"`
	Dosen   Dosen  `gorm:"foreignKey:NIP;references:NIP"`
	Kelas   Kelas  `gorm:"foreignKey:IDKelas"`
}

func (DosenPengampu) TableName() string {
	return "dosen_pengampu"
}
