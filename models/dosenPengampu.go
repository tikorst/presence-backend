package models

type DosenPengampu struct {
	NIP     string  `gorm:"primaryKey" grom:"unique" json:"nip"`
	IDKelas int     `gorm:"primaryKey" json:"id_kelas"`
	Status  string  `json:"status"`
	Kelas   []Kelas `gorm:"foreignKey:IDKelas;references:IDKelas"`
}

func (DosenPengampu) TableName() string {
	return "dosen_pengampu"
}
