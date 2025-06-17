package models

import "time"

type UserLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id_user_log"`
	IDUser    uint      `json:"id_user"`
	DeviceID  string    `json:"device_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	LoginTime time.Time `json:"login_time"`
	Success   bool      `json:"success"`
}
