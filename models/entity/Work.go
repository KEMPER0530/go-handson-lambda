package entity

import "time"

type Work struct {
	Work_id    int       `gorm:"primary_key;not null"    json:"work_id"`
	Title      string    `gorm:"type:varchar(256);null"  json:"title"`
	Comment    string    `gorm:"type:varchar(2000);null" json:"comment"`
	Img        string    `gorm:"type:varchar(256);null"  json:"img"`
	Url        string    `gorm:"type:varchar(256);null"  json:"url"`
	Ref        string    `gorm:"type:varchar(256);null"  json:"ref"`
	Updated_at time.Time `gorm:"type:null"         json:"updated_at"`
}
