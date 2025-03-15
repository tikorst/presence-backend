package models

type MataKuliah struct {
	KodeMatkul   string `gorm:"primaryKey" json:"kode_matkul"`
	NamaMatkul   string `json:"nama_matkul"`
	NamaMatkulEn string `json:"nama_matkul_en"`
	Prasyarat    string `json:"prasyarat"`
	Deskripsi    string `json:"deskripsi" gorm:"type:text"`
	SKSTeori     *int   `json:"sks_teori"`
	SKSPraktikum *int   `json:"sks_praktikum"`
	SKSPraktik   *int   `json:"sks_praktik"`
	TotalSKS     *int   `json:"total_sks"`
	Status       string `json:"status"`
}

func (MataKuliah) TableName() string {
	return "mata_kuliah"
}
