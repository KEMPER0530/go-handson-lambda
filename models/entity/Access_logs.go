package entity

import "time"

type Access_logs struct {
	User_id        string    `gorm:"type:varchar(256);null" json:"user_id"`
	Event_id       string    `gorm:"type:varchar(256);null" json:"event_id"`
	Access_ip      string    `gorm:"type:varchar(256);null" json:"access_ip"`
	City           string    `gorm:"type:varchar(256);null" json:"city"`
	Region         string    `gorm:"type:varchar(256);null" json:"region"`
	Region_code    string    `gorm:"type:varchar(8);null"   json:"region_code"`
	Country_name   string    `gorm:"type:varchar(64);null"  json:"country_name"`
	Country_code   string    `gorm:"type:varchar(8);null"   json:"country_code"`
	Continent_name string    `gorm:"type:varchar(64);null"  json:"continent_name"`
	Continent_code string    `gorm:"type:varchar(8);null"   json:"continent_code"`
	Latitude       string    `gorm:"type:varchar(128);null" json:"latitude"`
	Longitude      string    `gorm:"type:varchar(128);null" json:"longitude"`
	Postal         string    `gorm:"type:varchar(64);null"  json:"postal"`
	Calling_code   string    `gorm:"type:varchar(8);null"   json:"calling_code"`
	Created_at     time.Time `gorm:"type:null"              json:"created_at"`
}
