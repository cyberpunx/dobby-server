package model

import (
	"time"
)

type UserSession struct {
	IsLoggedIn              bool       `json:"isLoggedIn"`
	Username                *string    `json:"username"`
	Initials                *string    `json:"initials"`
	LoginDatetime           *time.Time `json:"datetime"`
	UserDateFormat          *string    `json:"userDateFormat"`
	IsCorrectDateFmt        bool       `json:"isCorrectDateFmt"`
	IsCorrectTimeZone       bool       `json:"isCorrectTimeZone"`
	IsCorrectTimeFmtAndZone bool       `json:"isCorrectTimeFmtAndZone"`
	User                    *User
	Permissions             []Permission
	PostSecret1             *string
	PostSecret2             *string
	ForumCookies            []CookieEntry
	UserDateTime            *time.Time
}

type CookieEntry struct {
	Name       string
	Value      string
	Domain     string
	Path       string
	SameSite   string
	Secure     bool
	HttpOnly   bool
	Persistent bool
	HostOnly   bool
	Expires    time.Time
	Creation   time.Time
	LastAccess time.Time

	// seqNum is a sequence number so that Cookies returns cookies in a
	// deterministic order, even for cookies that have equal Path length and
	// equal Creation time. This simplifies testing.
	seqNum uint64
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
