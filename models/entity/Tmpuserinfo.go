package entity

import "time"

type Tmpuserinfo struct {
	Email      string    `gorm:"type:varchar(20);not null" json:"email"`
	Password   string    `gorm:"type:varchar(200);not null" json:"password"`
	Name       string    `gorm:"type:varchar(200);not null" json:"name"`
	Token      string    `gorm:"type:varchar(256);not null" json:"token"`
	Expired    time.Time `gorm:"type:null"                  json:"expired"`
	Updated_at time.Time `gorm:"type:null"                  json:"updated_at"`
}
