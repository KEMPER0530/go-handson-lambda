package entity

type Mail_send_inf struct {
	Msg_id        string `gorm:"type:char(36);not null" json:"msg_id"`
	From_email    string `gorm:"type:varchar(256);null" json:"from_email"`
	To_email      string `gorm:"type:varchar(256);null" json:"to_email"`
	Subject       string `gorm:"type:varchar(100);null" json:"subject"`
	Body          string `gorm:"type:text;null"         json:"body"`
	Personal_name string `gorm:"type:varchar(64);null"  json:"personal_name"`
}
