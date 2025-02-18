package models

import "time"

type Session struct {
	Token  string    `gorm:"type:varchar(43);primaryKey;"`
	Data   []byte    `gorm:"type:bytea;notnull;"`
	Expiry time.Time `gorm:"type:timestamp(6);notnull;index;"`
}

func(Session) TableName() string {
	return "sessions"
}
