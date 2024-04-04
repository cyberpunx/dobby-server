package tursodb

import "strings"

func NoSuchTableError(err error) bool {
	return strings.Contains(err.Error(), NoSuchTable)
}
