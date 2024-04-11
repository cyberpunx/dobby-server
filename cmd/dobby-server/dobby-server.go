package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"localdev/dobby-server/internal/app/dobby-server/handler"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/storage"
	mylogger "localdev/dobby-server/internal/pkg/log"
	"localdev/dobby-server/internal/pkg/util"
	"os"
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
	// Gets environment variables
	tursoDbUrl := os.Getenv("TURSO_DB_URL")
	tursoDbToken := os.Getenv("TURSO_DB_TOKEN")
	serverPort := os.Getenv("SERVER_PORT")

	// Connects to TursoDB and gets config table
	tursoConnectLine := fmt.Sprintf(tursoDbUrl + "?authToken=" + tursoDbToken)
	store := storage.NewStore(tursoConnectLine)
	configApi := model.NewConfigApi(model.Config{}, *store)
	configTable, err := configApi.GetConfig()
	util.Panic(err)

	fmt.Printf("Config: \n %s", util.MarshalJsonPretty(configTable))

	// Starts the server
	app := echo.New()
	app.Use(slogecho.New(mylogger.GetLogger()))

	app.Static("/assets", "/internal/app/dobby-server/assets")

	handler.SetupRoutes(app, &configTable, store)

	err = app.Start(":" + serverPort)
	util.Panic(err)
}
