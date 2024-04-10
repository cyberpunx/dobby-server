package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"localdev/dobby-server/internal/app/dobby-server/config"
	"localdev/dobby-server/internal/app/dobby-server/handler"
	"localdev/dobby-server/internal/app/dobby-server/storage/tursodb"
	mylogger "localdev/dobby-server/internal/pkg/log"
	"localdev/dobby-server/internal/pkg/util"
)

const (
	Port = "8080"
)

func main() {
	//Loads .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found, not loaded")
	}

	// Loads config file
	configFile := config.GetConfigFile()
	configFile.Validate()

	// Connects to TursoDB and gets config table
	tursoConnectLine := fmt.Sprintf(*configFile.TursoDbUrl + "?authToken=" + *configFile.TursoDbToken)
	store := tursodb.InitDB(tursoConnectLine)
	configTable := store.GetConfig()

	// Merges config file and config table
	conf := config.MergeConfigs(*configFile, *configTable)
	fmt.Printf("Config: \n %s", util.MarshalJsonPretty(conf))

	// Starts the server
	app := echo.New()
	app.Use(slogecho.New(mylogger.GetLogger()))

	app.Static("/assets", "/internal/app/dobby-server/assets")

	handler.SetupRoutes(app, &conf, store)

	err = app.Start(":" + *conf.ServerPort)
	util.Panic(err)
}
