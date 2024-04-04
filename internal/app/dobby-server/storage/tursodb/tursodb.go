package tursodb

import (
	"database/sql"
	"fmt"
	"localdev/dobby-server/internal/app/dobby-server/config"
	"localdev/dobby-server/internal/app/dobby-server/storage/query"
	"localdev/dobby-server/internal/pkg/util"
	"os"
)

const (
	NoSuchTable    = "no such table"
	ForumBaseUrl   = "https://www.hogwartsrol.com/"
	PotionsClubUrl = "f98-club-de-pociones"
)

type Store struct {
	Conn *sql.DB
}

func InitDB(connectionUrl string) *Store {
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
			_, err = db.Exec(query.InsertConfigTable, ForumBaseUrl, "token.json", "client_secret.json", "13CCYZ4veljB6ItPNHdvxvClBZJaC1w-QMkq-H5btR74")
			util.Panic(err)
		} else {
			util.Panic(err)
		}
	} else {
		util.Panic(err)
	}

	_, err = db.Query(query.SelectPotionSubTable)
	if err != nil {
		if NoSuchTableError(err) {
			_, err = db.Exec(query.CreatePotionSubTable)
			util.Panic(err)
			_, err = db.Exec(query.InsertPotionSubTable, PotionsClubUrl, 72, 8)
			util.Panic(err)
		} else {
			util.Panic(err)
		}
	} else {
		util.Panic(err)
	}

	_, err = db.Query(query.SelectPotionThrTable)
	if err != nil {
		if NoSuchTableError(err) {
			_, err = db.Exec(query.CreatePotionThrTable)
			util.Panic(err)
		} else {
			util.Panic(err)
		}
	} else {
		util.Panic(err)
	}

	return &Store{Conn: db}
}

func (store Store) GetConfig() *config.Config {
	var conf config.Config
	db := store.Conn

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

func (store Store) UpdateConfig(config *config.Config) {
	db := store.Conn
	_, err := db.Exec(query.EnsureConfigRow)
	util.Panic(err)

	_, err = db.Exec(query.UpdateConfigTable,
		config.BaseUrl,
		config.GSheetTokenFile,
		config.GSheetCredFile,
		config.GSheetId)
	util.Panic(err)
}

func (store Store) GetPotionThread() *[]config.PotionThreadConfig {
	db := store.Conn
	rows, err := db.Query(query.SelectPotionThrTable)
	util.Panic(err)

	var threads []config.PotionThreadConfig
	for rows.Next() {
		var thread config.PotionThreadConfig
		err = rows.Scan(
			&thread.Url,
			&thread.TimeLimit,
			&thread.TurnLimit)
		util.Panic(err)
		threads = append(threads, thread)
	}

	return &threads
}

func (store Store) UpdatePotionThread(potionThrConfig *[]config.PotionThreadConfig) {
	db := store.Conn
	// Truncate the table and insert the new values
	_, err := db.Exec(query.TruncateTable, "PotionThreadConfig")
	util.Panic(err)

	//insert one by one the potionThrConfig
	for _, potionThread := range *potionThrConfig {
		_, err := db.Exec(query.InsertPotionThrTable, potionThread.Url, potionThread.TimeLimit, potionThread.TurnLimit)
		util.Panic(err)
	}
}

func (store Store) GetPotionSubforum() *[]config.PotionSubforumConfig {
	db := store.Conn
	rows, err := db.Query(query.SelectPotionSubTable)
	util.Panic(err)

	var subforums []config.PotionSubforumConfig
	for rows.Next() {
		var subforum config.PotionSubforumConfig
		err = rows.Scan(
			&subforum.Url,
			&subforum.TimeLimit,
			&subforum.TurnLimit)
		util.Panic(err)
		subforums = append(subforums, subforum)
	}

	return &subforums
}

func (store Store) UpdatePotionSubforum(potionSubConfig *[]config.PotionSubforumConfig) {
	db := store.Conn
	// Truncate the table and insert the new values
	_, err := db.Exec(query.TruncateTable, "PotionSubforumConfig")
	util.Panic(err)

	//insert one by one the potionSubConfig
	for _, potionSubforum := range *potionSubConfig {
		_, err := db.Exec(query.InsertPotionSubTable, potionSubforum.Url, potionSubforum.TimeLimit, potionSubforum.TurnLimit)
		util.Panic(err)
	}
}
