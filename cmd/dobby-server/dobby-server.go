package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"localdev/dobby-server/internal/app/dobby-server/config"
	"localdev/dobby-server/internal/app/dobby-server/handler"
	"localdev/dobby-server/internal/app/dobby-server/storage/tursodb"
	"localdev/dobby-server/internal/pkg/util"
)

const (
	Port = "8080"
)

func main() {
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
	handler.SetupRoutes(app, &conf, store)

	err := app.Start(":" + *conf.ServerPort)
	util.Panic(err)
}
