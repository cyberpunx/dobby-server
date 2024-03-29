package config

import (
	"encoding/json"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetConfigFile() *Config {
	return config
}

func LoadConfigFile(configPath string, config interface{}) (string, error) {
	executablePath, err := os.Executable()
	parentPath := filepath.Dir(executablePath)
	//abs, err := filepath.Abs(path)
	//util.Panic(err)
	loadedPath := filepath.Join(parentPath, configPath)
	b, err := ioutil.ReadFile(filepath.Join(parentPath, configPath))
	if err != nil {
		loadedPath = configPath
		b, err = ioutil.ReadFile(filepath.Join(configPath))
		if err != nil {
			return "", err
		}
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		return "", err
	}
	return loadedPath, nil
}

func InitUnicodeConfig(conf *Config) {
	if *conf.UnicodeOutput {
		Reset = "\033[0m"
		Red = "\033[31m"
		Green = "\033[32m"
		Yellow = "\033[33m"
		Blue = "\033[34m"
		Purple = "\033[35m"
		Cyan = "\033[36m"
		Gray = "\033[37m"
		White = "\033[97m"
		CheckEmoji = "✔"
		CrossEmoji = "❌"
		RightArrowEmoji = "▶"
	} else {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
		CheckEmoji = "[OK]"
		CrossEmoji = "[X]"
		RightArrowEmoji = "-->"
	}
}

func init() {
	conf := flag.String("conf", "./config/conf.json", "Config")
	_, err := LoadConfigFile(*conf, &config)
	if err != nil {
		//if cannot load config file, try get Environment variables
		tursoDbUrl := os.Getenv("TURSO_DB_URL")
		tursoDbToken := os.Getenv("TURSO_DB_TOKEN")
		serverPort := os.Getenv("SERVER_PORT")
		unicodeOutput := true
		if tursoDbUrl == "" || tursoDbToken == "" || serverPort == "" {
			panic("Cannot load config file and Environment variables are not set")

		}
		config = &Config{
			TursoDbUrl:    &tursoDbUrl,
			TursoDbToken:  &tursoDbToken,
			ServerPort:    &serverPort,
			UnicodeOutput: &unicodeOutput,
		}

	}
	InitUnicodeConfig(config)
}
