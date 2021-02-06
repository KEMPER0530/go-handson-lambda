package entity

import "time"

type Mst_ssmlknr struct {
	Id            int       `gorm:"primary_key" json:"id"`
	Subject       string    `gorm:"type:varchar(100);not null" json:"subject"`
	Body          string    `gorm:"type:varchar(2000);not null" json:"body"`
	Replytitle    string    `gorm:"type:varchar(100);not null" json:"replytitle"`
	Replyconcents string    `gorm:"type:char(2000);null" json:"replyconcents"`
	Toreply       string    `gorm:"type:varchar(100);null" json:"toreply"`
	Updated_at    time.Time `gorm:"type:null" json:"updated_at"`
}
