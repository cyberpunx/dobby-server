package model

import "localdev/dobby-server/internal/app/dobby-server/storage"

const (
	CreateAnnouncementTable = `CREATE TABLE IF NOT EXISTS Announcement (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			"title" TEXT NOT NULL,
    			"message" TEXT NOT NULL
	);`

	SelectAnnouncementTable = `SELECT * FROM Announcement;`
	InsertAnnouncementTable = `INSERT INTO Announcement (
                          				title,
                          				message)
    					  				VALUES (?, ?);`
	UpdateAnnouncementTable = `UPDATE Announcement SET
						title = ?,
						message = ?
						WHERE id = ?;`
	DeleteAnnouncementTable = `DELETE FROM Announcement WHERE id = ?;`
)

type Announcement struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
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
		err = rows.Scan(&a.Id, &a.Title, &a.Message)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}
	return announcements, nil
}

func (api *AnnouncementApi) InsertAnnouncement(a *Announcement) error {
	_, err := api.Store.Conn.Exec(DeleteAnnouncementTable, a.Title, a.Message)
	if err != nil {
		return err
	}
	return nil
}

func (api *AnnouncementApi) UpdateAnnouncement(id int, a *Announcement) error {
	_, err := api.Store.Conn.Exec(UpdateAnnouncementTable, a.Title, a.Message, id)
	if err != nil {
		return err
	}
	return nil
}

func (api *AnnouncementApi) DeleteUserById(id int) error {
	_, err := api.Store.Conn.Exec(DeleteUserTable, id)
	if err != nil {
		return err
	}
	return nil
}
