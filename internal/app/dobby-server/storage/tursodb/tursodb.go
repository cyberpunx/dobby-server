package tursodb

import (
	"database/sql"
	"fmt"
	"localdev/dobby-server/internal/app/dobby-server/config"
	"localdev/dobby-server/internal/app/dobby-server/storage/query"
	"localdev/dobby-server/internal/pkg/util"
	"os"
)

func InitDB(connectionUrl string) *sql.DB {
	db, err := sql.Open("libsql", connectionUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", connectionUrl, err)
		os.Exit(1)
	}

	_, err = db.Query(query.SelectConfigTable)
	if err != nil {
		if NoSuchTableError(err) {
			_, err = db.Exec(query.CreateConfigTable)
			util.Panic(err)
			_, err = db.Exec(query.InsertConfigTable, "https://www.hogwartsrol.com/", "token.json", "client_secret.json", "13CCYZ4veljB6ItPNHdvxvClBZJaC1w-QMkq-H5btR74")
			util.Panic(err)
		} else {
			util.Panic(err)
		}
	} else {
		util.Panic(err)
	}

	return db
}

func GetConfig(db *sql.DB) *config.Config {
	var conf config.Config

	// Execute the query
	row := db.QueryRow(query.SelectConfigTable)
	err := row.Scan(
		&conf.BaseUrl,
		&conf.GSheetTokenFile,
		&conf.GSheetCredFile,
		&conf.GSheetId)
	util.Panic(err)

	return &conf
}
