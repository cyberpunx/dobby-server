package model

import (
	"errors"
	"localdev/dobby-server/internal/app/dobby-server/storage"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics"
	"strings"
)

const (
	SelectUserTable      = `SELECT * FROM User;`
	SelectUserByUsername = `SELECT * FROM User WHERE username = ?;`
	CreateUserTable      = `CREATE TABLE IF NOT EXISTS User (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	"username" TEXT UNIQUE,
    	"active" BOOLEAN NOT NULL,
    	"title" TEXT NOT NULL,
    	"permissions" TEXT NOT NULL
	);`
	InsertUserTable = `INSERT INTO User (
                  		username,
                  		active,
                  		title,
                  		permissions)
    			  		VALUES (?, ?, ?, ?);`
	UpdateUserTable = `UPDATE User SET
					username = ?,
					active = ?,
					title = ?,
					permissions = ?
					WHERE id = ?;`
	DeleteUserTable = `DELETE FROM User WHERE id = ?;`

	PermissionAdmin           Permission = "Admin"
	PermissionAdminReadOnly   Permission = "AdminReadOnly"
	PermissionPotions         Permission = Permission(dynamics.DynamicPotion)
	PermissionCreationChamber Permission = Permission(dynamics.DynamicCreationChamber)
)

type Permission string

type User struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Active      bool   `json:"active"`
	Title       string `json:"title"`
	Permissions string `json:"permissions"`
}

type UserCrud struct {
	User
	EditUrl   string
	ViewUrl   string
	UpdateUrl string
	DeleteUrl string
}

type UserApi struct {
	User  User
	Store storage.Store
}

func GetAllPermissions() []Permission {
	return []Permission{
		PermissionAdmin,
		PermissionAdminReadOnly,
		PermissionPotions,
		PermissionCreationChamber,
	}
}

func NewUserApi(u User, store storage.Store) *UserApi {
	return &UserApi{
		User:  u,
		Store: store,
	}
}

func (api *UserApi) CreateInitialUserTable() error {
	_, err := api.Store.Conn.Exec(CreateUserTable)
	if err != nil {
		return err
	}
	return nil
}

func (api *UserApi) GetAllUser() ([]User, error) {
	rows, err := api.Store.Conn.Query(SelectUserTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		err = rows.Scan(&api.User.Id, &api.User.Username, &api.User.Active, &api.User.Title, &api.User.Permissions)
		users = append(users, api.User)
	}
	return users, nil
}

func (api *UserApi) InsertUser(user User) error {
	_, err := api.Store.Conn.Exec(InsertUserTable, user.Username, user.Active, user.Title, user.Permissions)
	if err != nil {
		return err
	}
	return nil
}

func (api *UserApi) UpdateUser(id int, user *User) error {
	_, err := api.Store.Conn.Exec(UpdateUserTable, user.Username, user.Active, user.Title, user.Permissions, id)
	if err != nil {
		return err
	}
	return nil
}

func (api *UserApi) GetUserByUsername(username string) (*User, error) {
	rows, err := api.Store.Conn.Query(SelectUserByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user User
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Active, &user.Title, &user.Permissions)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("no user found")
	}
	if rows.Next() { // Check if there is more than one row
		return nil, errors.New("multiple users found with the same username")
	}
	return &user, nil
}

func (api *UserApi) DeleteUserById(id int) error {
	_, err := api.Store.Conn.Exec(DeleteUserTable, id)
	if err != nil {
		return err
	}
	return nil
}

func (api *UserApi) GetUserById(id int) (*User, error) {
	rows, err := api.Store.Conn.Query(SelectUserTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user User
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Active, &user.Title, &user.Permissions)
		if err != nil {
			return nil, err
		}
		if user.Id == id {
			return &user, nil
		}
	}
	return nil, errors.New("no user found")
}

func PermissionsToString(permissions []string) string {
	if len(permissions) == 0 {
		return ""
	}
	var permissionsString string
	for _, permission := range permissions {
		permissionsString += permission + ","
	}
	return permissionsString[:len(permissionsString)-1]
}

func StringToPermissions(permissionsString string) []string {
	if permissionsString == "" {
		return []string{}
	}
	return strings.Split(permissionsString, ",")
}

func (user *User) GetUserPermissions() []Permission {
	var permissions []Permission
	for _, permission := range StringToPermissions(user.Permissions) {
		permissions = append(permissions, Permission(permission))
	}
	return permissions
}

func (user *User) HavePermission(p Permission) bool {
	for _, perm := range user.GetUserPermissions() {
		if perm == p {
			return true
		}
	}
	return false
}
