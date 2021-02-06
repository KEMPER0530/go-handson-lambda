package entity

import "time"

type Crdcardinfo struct {
	Cardnumber string    `gorm:"type:char(19);not null" json:"cardnumber"`
	Cardname   string    `gorm:"type:varchar(60);null"  json:"cardname"`
	Cardmonth  int       `gorm:"null"                   json:"cardmonth"`
	Cardyear   int       `gorm:"null"                   json:"cardyear"`
	Cardcvv    string    `gorm:"type:varchar(256);null" json:"cardcvv"`
	Updated_at time.Time `gorm:"type:null"              json:"updated_at"`
}
