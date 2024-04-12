package model

import (
	"time"
)

type UserSession struct {
	IsLoggedIn       bool       `json:"isLoggedIn"`
	Username         *string    `json:"username"`
	Initials         *string    `json:"initials"`
	LoginDatetime    *time.Time `json:"datetime"`
	UserDateFormat   *string    `json:"userDateFormat"`
	IsCorrectDateFmt bool       `json:"isCorrectDateFmt"`
}
