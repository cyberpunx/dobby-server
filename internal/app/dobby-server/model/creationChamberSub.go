package model

import "localdev/dobby-server/internal/app/dobby-server/storage"

const (
	CreationChamberUrl       = "f237-camara-de-la-creacion-magica"
	CreationChamberTimeLimit = 72
	CreationChamberTurnLimit = 8

	SelectCreationChamberSubTable = `SELECT * FROM CreationChamberSubforumConfig;`
	CreateCreationChamberSubTable = `CREATE TABLE IF NOT EXISTS CreationChamberSubforumConfig (
        "url" TEXT PRIMARY KEY,
        "timeLimit" INTEGER NOT NULL,
        "turnLimit" INTEGER NOT NULL
    );`
	InsertCreationChamberSubTable = `INSERT INTO CreationChamberSubforumConfig (
		url, 
		timeLimit, 
		turnLimit)
		VALUES (?, ?, ?);`
)

type CreationChamberSub struct {
	Url       string `json:"url"`
	TimeLimit int    `json:"timeLimit"`
	TurnLimit int    `json:"turnLimit"`
}

type CreationChamberSubApi struct {
	CreationChamberSub CreationChamberSub
	Store              storage.Store
}

func NewCreationChamberSubApi(p CreationChamberSub, store storage.Store) *CreationChamberSubApi {
	return &CreationChamberSubApi{
		CreationChamberSub: p,
		Store:              store,
	}
}

func (api *CreationChamberSubApi) CreateInitialCreationChamberSubTable() (CreationChamberSub, error) {
	_, err := api.Store.Conn.Exec(CreateCreationChamberSubTable)
	if err != nil {
		return CreationChamberSub{}, err
	}
	_, err = api.Store.Conn.Exec(InsertCreationChamberSubTable, CreationChamberUrl, CreationChamberTimeLimit, CreationChamberTurnLimit)
	if err != nil {
		return CreationChamberSub{}, err
	}
	return CreationChamberSub{
		Url:       CreationChamberUrl,
		TimeLimit: 72,
		TurnLimit: 8,
	}, nil
}

func (api *CreationChamberSubApi) GetAllCreationChamberSub() ([]CreationChamberSub, error) {
	rows, err := api.Store.Conn.Query(SelectCreationChamberSubTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var potionSubs []CreationChamberSub
	for rows.Next() {
		err = rows.Scan(&api.CreationChamberSub.Url, &api.CreationChamberSub.TimeLimit, &api.CreationChamberSub.TurnLimit)
		potionSubs = append(potionSubs, api.CreationChamberSub)
	}
	return potionSubs, nil
}

func (api *CreationChamberSubApi) InsertCreationChamberSub(potionSubConfig *CreationChamberSub) error {
	_, err := api.Store.Conn.Exec(InsertCreationChamberSubTable, potionSubConfig.Url, potionSubConfig.TimeLimit, potionSubConfig.TurnLimit)
	if err != nil {
		return err
	}
	return nil
}
