package gsheet

import "time"

type DayOff struct {
	Username string
	Date     time.Time
}

type PlayerBonus struct {
	Username string
	Bonus    int
}
