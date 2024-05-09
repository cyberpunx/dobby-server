package model

import (
	"errors"
	"localdev/dobby-server/internal/app/dobby-server/storage"
)

const (
	CreateAnnouncementTable = `CREATE TABLE IF NOT EXISTS Announcement (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			"title" TEXT NOT NULL,
    			"message" TEXT NOT NULL,
				"type" TEXT NOT NULL
	);`

	SelectAnnouncementTable = `SELECT * FROM Announcement;`
	SelectAnnouncementById  = `SELECT * FROM Announcement WHERE id = ?;`
	InsertAnnouncementTable = `INSERT INTO Announcement (
                          				title,
                          				message,
                          				type)
    					  				VALUES (?, ?, ?);`
	UpdateAnnouncementTable = `UPDATE Announcement SET
						title = ?,
						message = ?,
						type = ?
						WHERE id = ?;`
	DeleteAnnouncementTable = `DELETE FROM Announcement WHERE id = ?;`

	TypeGeneral    AnnouncementType = "General"
	TypeKnownIssue AnnouncementType = "Known Issue"
)

func GetAllAnnouncementTypes() []AnnouncementType {
	return []AnnouncementType{
		TypeGeneral,
		TypeKnownIssue,
	}
}

type AnnouncementType string

type Announcement struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type AnnouncementCrud struct {
	Announcement
	EditUrl   string
	ViewUrl   string
	UpdateUrl string
	DeleteUrl string
}

type AnnouncementApi struct {
	Announcement Announcement
	Store        storage.Store
}

func NewAnnouncementApi(a Announcement, store storage.Store) *AnnouncementApi {
	return &AnnouncementApi{
		Announcement: a,
		Store:        store,
	}
}

func (api *AnnouncementApi) CreateAnnouncementTable() error {
	_, err := api.Store.Conn.Exec(CreateAnnouncementTable)
	if err != nil {
		return err
	}
	return nil
}

func (api *AnnouncementApi) GetAllAnnouncement() ([]Announcement, error) {
	rows, err := api.Store.Conn.Query(SelectAnnouncementTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []Announcement
	for rows.Next() {
		var a Announcement
		err = rows.Scan(&a.Id, &a.Title, &a.Message, &a.Type)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}
	return announcements, nil
}

func (api *AnnouncementApi) InsertAnnouncement(a *Announcement) error {
	_, err := api.Store.Conn.Exec(InsertAnnouncementTable, a.Title, a.Message, a.Type)
	if err != nil {
		return err
	}
	return nil
}

func (api *AnnouncementApi) UpdateAnnouncement(id int, a *Announcement) error {
	_, err := api.Store.Conn.Exec(UpdateAnnouncementTable, a.Title, a.Message, a.Type, id)
	if err != nil {
		return err
	}
	return nil
}

func (api *AnnouncementApi) DeleteAnnouncementById(id int) error {
	_, err := api.Store.Conn.Exec(DeleteAnnouncementTable, id)
	if err != nil {
		return err
	}
	return nil
}

func (api *AnnouncementApi) GetAnnouncementById(id int) (*Announcement, error) {
	rows, err := api.Store.Conn.Query(SelectAnnouncementById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var announcement Announcement
	for rows.Next() {
		err = rows.Scan(&announcement.Id, &announcement.Title, &announcement.Message, &announcement.Type)
		if err != nil {
			return nil, err
		}
		if announcement.Id == id {
			return &announcement, nil
		}
	}
	return nil, errors.New("no user found")
}
