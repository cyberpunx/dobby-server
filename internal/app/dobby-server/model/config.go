package model

import (
	"localdev/dobby-server/internal/app/dobby-server/storage"
	"localdev/dobby-server/internal/pkg/util"
)

const (
	ForumBaseUrl = "https://www.hogwartsrol.com/"

	SelectConfigTable       = `SELECT * FROM Config;`
	SelectConfigTableLimit1 = `SELECT * FROM Config LIMIT 1;`
	CreateConfigTable       = `CREATE TABLE IF NOT EXISTS Config (
        "baseUrl" TEXT,
        "gSheetTokenFile" TEXT,
        "gSheetCredFile" TEXT
    );`
	InsertConfigTable = `INSERT INTO Config (
		baseUrl,
		gSheetTokenFile,
		gSheetCredFile)
		VALUES (?, ?, ?);`
	UpdateConfigTable = `UPDATE Config SET 
		baseUrl = ?, 
		gSheetTokenFile = ?, 
		gSheetCredFile = ?;`
	EnsureConfigRow = `INSERT INTO Config (
					baseUrl, 
					gSheetTokenFile, 
					gSheetCredFile)
					SELECT '', '', '', '' WHERE NOT EXISTS (SELECT 1 FROM Config);`
)

type Config struct {
	BaseUrl         string `json:"baseUrl"`
	GSheetTokenFile string `json:"gSheetTokenFile"`
	GSheetCredFile  string `json:"gSheetCredFile"`
}

type ConfigApi struct {
	Config Config
	Store  storage.Store
}

func NewConfigApi(c Config, store storage.Store) *ConfigApi {
	return &ConfigApi{
		Config: c,
		Store:  store,
	}
}

func (api *ConfigApi) CreateInitialConfigTable() error {
	_, err := api.Store.Conn.Exec(CreateConfigTable)
	if err != nil {
		return err
	}

	configTable, err := api.GetAllConfig()
	if len(configTable) == 0 {
		_, err = api.Store.Conn.Exec(InsertConfigTable, ForumBaseUrl, "token.json", "client_secret.json")
		if err != nil {
			return err
		}

	}
	return nil
}

func (api *ConfigApi) GetConfig() (Config, error) {
	rows, err := api.Store.Conn.Query(SelectConfigTableLimit1)
	util.Panic(err)
	defer rows.Close()
	var c Config
	for rows.Next() {
		err = rows.Scan(&c.BaseUrl, &c.GSheetTokenFile, &c.GSheetCredFile)
	}
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func (api *ConfigApi) InsertConfig(c Config) (Config, error) {
	_, err := api.Store.Conn.Exec(UpdateConfigTable, c.BaseUrl, c.GSheetTokenFile, c.GSheetCredFile)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func (api *ConfigApi) GetAllConfig() ([]Config, error) {
	rows, err := api.Store.Conn.Query(SelectConfigTable)
	util.Panic(err)
	defer rows.Close()
	var configs []Config
	for rows.Next() {
		err = rows.Scan(&api.Config.BaseUrl, &api.Config.GSheetTokenFile, &api.Config.GSheetCredFile)
		configs = append(configs, api.Config)
	}
	return configs, nil
}
