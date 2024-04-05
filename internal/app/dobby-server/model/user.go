package model

import "time"

type User struct {
	IsLoggedIn bool       `json:"isLoggedIn"`
	Username   *string    `json:"username"`
	Initials   *string    `json:"initials"`
	Datetime   *time.Time `json:"datetime"`
}
