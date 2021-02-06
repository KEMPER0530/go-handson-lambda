package entity

import "time"

type ActuatorRslt struct {
	App_Status string
	Db_Status  string
	Time       time.Time
	Host       string
}
