package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"localdev/dobby-server/internal/pkg/util"
	"path/filepath"
)

const (
	GobsPostDynamic  PostDynamic = "GOBS"
	ChessPostDynamic PostDynamic = "CHESS"
)

type PostDynamic string

type Config struct {
	Users         []User `json:"users"`
	BaseUrl       string `json:"baseUrl"`
	GobstonsUrl   string `json:"gobstonsUrl"`
	MagicChessUrl string `json:"magicChessUrl"`
}

type User struct {
	Username string        `json:"username"`
	Password string        `json:"password"`
	Msg      []UserMessage `json:"msg"`
}

type UserMessage struct {
	Message     string      `json:"message"`
	Subject     string      `json:"subject"`
	SubForumId  string      `json:"subForumId"`
	SubForumUrl string      `json:"subForumUrl"`
	PostDynamic PostDynamic `json:"postDynamic"`
}

func ReadConfigFile(path string) Config {
	var config Config
	abs, err := filepath.Abs(path)
	fmt.Println("Path: ", abs)
	util.Panic(err)
	b, err := ioutil.ReadFile(abs)
	util.Panic(err)
	err = json.Unmarshal(b, &config)
	util.Panic(err)
	return config
}

func NewConfig() Config {
	var users []User
	users = append(users, User{Username: "", Password: ""})
	return Config{
		Users:         users,
		BaseUrl:       "https://www.hogwartsrol.com/",
		GobstonsUrl:   "f34-gobstons",
		MagicChessUrl: "f177-ajedrez-magico",
	}
}
