package model

import "localdev/dobby-server/internal/app/dobby-server/storage"

const (
	PotionsClubUrl       = "f98-club-de-pociones"
	PotionsClubTimeLimit = 72
	PotionsClubTurnLimit = 8

	SelectPotionSubTable = `SELECT * FROM PotionSubforumConfig;`
	CreatePotionSubTable = `CREATE TABLE IF NOT EXISTS PotionSubforumConfig (
        "url" TEXT PRIMARY KEY,
        "timeLimit" INTEGER NOT NULL,
        "turnLimit" INTEGER NOT NULL
    );`
	InsertPotionSubTable = `INSERT INTO PotionSubforumConfig (
		url, 
		timeLimit, 
		turnLimit)
		VALUES (?, ?, ?);`
)

type PotionSub struct {
	Url       string `json:"url"`
	TimeLimit int    `json:"timeLimit"`
	TurnLimit int    `json:"turnLimit"`
}

type PotionSubApi struct {
	PotionSub PotionSub
	Store     storage.Store
}

func NewPotionSubApi(p PotionSub, store storage.Store) *PotionSubApi {
	return &PotionSubApi{
		PotionSub: p,
		Store:     store,
	}
}

func (api *PotionSubApi) CreateInitialPotionSubTable() (PotionSub, error) {
	_, err := api.Store.Conn.Exec(CreatePotionSubTable)
	if err != nil {
		return PotionSub{}, err
	}
	_, err = api.Store.Conn.Exec(InsertPotionSubTable, PotionsClubUrl, PotionsClubTimeLimit, PotionsClubTurnLimit)
	if err != nil {
		return PotionSub{}, err
	}
	return PotionSub{
		Url:       PotionsClubUrl,
		TimeLimit: 72,
		TurnLimit: 8,
	}, nil
}

func (api *PotionSubApi) GetAllPotionSub() ([]PotionSub, error) {
	rows, err := api.Store.Conn.Query(SelectPotionSubTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var potionSubs []PotionSub
	for rows.Next() {
		rows.Scan(&api.PotionSub.Url, &api.PotionSub.TimeLimit, &api.PotionSub.TurnLimit)
		potionSubs = append(potionSubs, api.PotionSub)
	}
	return potionSubs, nil
}

func (api *PotionSubApi) InsertPotionSub(p PotionSub) (PotionSub, error) {
	_, err := api.Store.Conn.Exec(InsertPotionSubTable, p.Url, p.TimeLimit, p.TurnLimit)
	if err != nil {
		return PotionSub{}, err
	}
	return PotionSub{
		Url:       p.Url,
		TimeLimit: p.TimeLimit,
		TurnLimit: p.TurnLimit,
	}, nil
}
