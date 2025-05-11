package models

import "time"

type User struct {
	IDUser            int       `gorm:"primaryKey;autoIncrement" json:"id_user"`
	Username          string    `gorm:"unique" json:"username"`
	Email             string    `json:"email"`
	Nama              string    `json:"nama"`
	Alamat            string    `json:"alamat"`
	TanggalLahir      time.Time `json:"tanggal_lahir"`
	TempatLahir       string    `json:"tempat_lahir"`
	JenisKelamin      string    `json:"jenis_kelamin"`
	NoTelepon         string    `json:"no_telepon"`
	Password          string    `json:"-"`
	Status            string    `json:"status"`
	TipeUser          string    `json:"tipe_user"`
	DeviceID          string    `json:"device_id" gorm:"default:NULL"`
	DeviceIDUpdatedAt *time.Time `json:"device_id_updated_at" gorm:"default:NULL"`
}
