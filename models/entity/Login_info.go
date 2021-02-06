package entity

import "time"

type Login_info struct {
	Username   string    `gorm:"type:varchar(20);not null" json:"username"`
	Password   string    `gorm:"type:varchar(200);not null" json:"password"`
	Name       string    `gorm:"type:varchar(200);not null" json:"name"`
	Id         int       `gorm:"null"                       json:"id"`
	Updated_at time.Time `gorm:"type:null"                  json:"updated_at"`
}
