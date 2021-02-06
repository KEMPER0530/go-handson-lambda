package entity

type Testmember struct {
	Id       int    `gorm:"primary_key;not null"     json:"id"`
	Name     string `gorm:"type:varchar(200);null"   json:"name"`
	Position string `gorm:"type:varchar(1000);null"  json:"position"`
	Height   int    `gorm:"null"                     json:"height"`
	Profile  string `gorm:"type:varchar(2000);null"  json:"profile"`
}
