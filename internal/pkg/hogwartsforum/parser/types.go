package parser

import "time"

type Thread struct {
	Title    string
	Url      string
	Author   *User
	Created  *time.Time
	LastPost *Post
	Pages    []string
	Posts    []*Post
}

type Post struct {
	Url     string
	Author  *User
	Created *time.Time
	Edited  *time.Time
	Content string
	Id      string
	Dices   []*Dice
}

type Dice struct {
	DiceLine string
	Result   int
}

type User struct {
	Username string
	Url      string
	House    string
}
