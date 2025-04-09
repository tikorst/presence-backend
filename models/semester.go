package models

type Semester struct {
	IDSemester  int    `gorm:"primaryKey;autoIncrement" json:"id_semester"`
	TahunAjaran string `json:"tahun_ajaran"`
	Semester    string `json:"semester"`
	Status      string `json:"status"`
}

func (Semester) TableName() string {
	return "semester"
}
