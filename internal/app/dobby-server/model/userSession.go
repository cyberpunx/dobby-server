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
	User             *User
	Permissions      []Permission
}

func (us *UserSession) IsUserLoggedIn() bool {
	return us.IsLoggedIn
}

func (us *UserSession) GetUserPermissions() []Permission {
	return us.Permissions
}

func (us *UserSession) HavePermission(p Permission) bool {
	for _, perm := range us.Permissions {
		if perm == p {
			return true
		}
	}
	return false
}
