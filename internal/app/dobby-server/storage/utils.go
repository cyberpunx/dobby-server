package storage

import "strings"

const (
	NoSuchTable = "no such table"
)

func NoSuchTableError(err error) bool {
	return strings.Contains(err.Error(), NoSuchTable)
}
