package models

type Nilai struct {
	IDNilai    int     `gorm:"primaryKey;autoIncrement" json:"id_nilai"`
	IDKelas    int     `json:"id_kelas"`
	NPM        string  `json:"npm"`
	NilaiUTS   float64 `json:"nilai_uts"`
	NilaiUAS   float64 `json:"nilai_uas"`
	Bobot      float64 `json:"bobot"`
	NilaiHuruf float64 `json:"nilai_huruf"`
}
