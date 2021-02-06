package entity

import "time"

type Mail_send_rslt struct {
	Send_no         int       `gorm:"primary_key" json:"send_no"`
	Msg_id          string    `gorm:"type:char(36);not null" json:"msg_id"`
	Msg_id_ses      string    `gorm:"type:char(256);null" json:"msg_id_ses"`
	Tnnt_id         string    `gorm:"type:char(5);not null" json:"tnnt_id"`
	Target_sys_type string    `gorm:"type:char(1);not null" json:"target_sys_type"`
	Status          string    `gorm:"type:char(1);null" json:"status"`
	Server_id       string    `gorm:"type:varchar(20);null" json:"server_id"`
	Queue_id        string    `gorm:"type:varchar(20);null" json:"queue_id"`
	Priority        int       `gorm:"null" json:"priority"`
	Send_reg_at     time.Time `gorm:"null" json:"send_reg_at"`
	Mta             string    `gorm:"type:varchar(5);null" json:"mta"`
	Dsn_cd          string    `gorm:"type:char(5);null" json:"dsn_cd"`
	Queue_remove    string    `gorm:"type:char(1);not null" json:"from_email"`
	Updated_at      time.Time `gorm:"null" json:"updated_at"`
}
