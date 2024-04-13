package storage

import (
	"database/sql"
	"localdev/dobby-server/internal/pkg/util"
	"os"
)

type Store struct {
	Conn *sql.DB
}

func NewStore(connectionUrl string) *Store {
	db, err := sql.Open("libsql", connectionUrl)
	if err != nil {
		util.LongPrintlnPrintln(os.Stderr, "failed to open db %s: %s", connectionUrl, err)
		os.Exit(1)
	}

	util.LongPrintlnPrintln("🚀 Connected successfully to the Database")
	return &Store{Conn: db}
}
