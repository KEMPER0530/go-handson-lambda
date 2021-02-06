package entity

import "time"

type Blklst_send_mailadr struct {
	Email           string `json:"email"`
	Updated_at      time.Time `gorm:"type:null"              json:"updated_at"`
}
