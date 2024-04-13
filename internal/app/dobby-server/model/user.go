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
    	"username" TEXT PRIMARY KEY,
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
					permissions = ?;`

	PermissionAdmin           Permission = "Admin"
	PermissionPotions         Permission = Permission(dynamics.DynamicPotion)
	PermissionCreationChamber Permission = Permission(dynamics.DynamicCreationChamber)
)

type Permission string

type User struct {
	Username    string `json:"username"`
	Active      bool   `json:"active"`
	Title       string `json:"title"`
	Permissions string `json:"permissions"`
}

type UserApi struct {
	User  User
	Store storage.Store
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
		err = rows.Scan(&api.User.Username, &api.User.Active, &api.User.Title, &api.User.Permissions)
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

func (api *UserApi) UpdateUser(user *User) error {
	_, err := api.Store.Conn.Exec(UpdateUserTable, user.Username, user.Active, user.Title, user.Permissions)
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
		err = rows.Scan(&user.Username, &user.Active, &user.Title, &user.Permissions)
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
