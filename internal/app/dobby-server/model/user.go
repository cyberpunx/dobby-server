package model

import "time"

type User struct {
	Username *string    `json:"username"`
	Initials *string    `json:"initials"`
	Datetime *time.Time `json:"datetime"`
}
