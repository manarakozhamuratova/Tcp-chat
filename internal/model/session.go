package model

import "time"

type Session struct {
	ID     int64
	User   User
	Token  string
	Expiry time.Time
}
