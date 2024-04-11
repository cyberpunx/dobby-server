package model

import "localdev/dobby-server/internal/app/dobby-server/storage"

const (
	SelectPotionThrTable = `SELECT * FROM PotionThreadConfig;`
	CreatePotionThrTable = `CREATE TABLE IF NOT EXISTS PotionThreadConfig (
        "url" TEXT PRIMARY KEY,
        "timeLimit" INTEGER NOT NULL,
        "turnLimit" INTEGER NOT NULL
    );`
	InsertPotionThrTable = `INSERT INTO PotionThreadConfig (
		url,
		timeLimit,
		turnLimit)
		VALUES (?, ?, ?);`
)

type PotionThread struct {
	Url       string `json:"url"`
	TimeLimit int    `json:"timeLimit"`
	TurnLimit int    `json:"turnLimit"`
}

type PotionThreadApi struct {
	PotionThread PotionThread
	Store        storage.Store
}

func NewPotionThreadApi(p PotionThread, store storage.Store) *PotionThreadApi {
	return &PotionThreadApi{
		PotionThread: p,
		Store:        store,
	}
}

func (api *PotionThreadApi) CreateInitialPotionThreadTable() error {
	_, err := api.Store.Conn.Exec(CreatePotionThrTable)
	if err != nil {
		return err
	}
	return nil
}

func (api *PotionThreadApi) GetAllPotionThread() ([]PotionThread, error) {
	rows, err := api.Store.Conn.Query(SelectPotionThrTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var potionThreads []PotionThread
	for rows.Next() {
		err = rows.Scan(&api.PotionThread.Url, &api.PotionThread.TimeLimit, &api.PotionThread.TurnLimit)
		potionThreads = append(potionThreads, api.PotionThread)
	}
	return potionThreads, nil
}

func (api *PotionThreadApi) InsertPotionThread(potionThrConfig *PotionThread) error {
	_, err := api.Store.Conn.Exec(InsertPotionThrTable, potionThrConfig.Url, potionThrConfig.TimeLimit, potionThrConfig.TurnLimit)
	if err != nil {
		return err
	}
	return nil
}
