package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"localdev/dobby-server/internal/app/dobby-server/handler"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/storage"
	mylogger "localdev/dobby-server/internal/pkg/logger"
	"localdev/dobby-server/internal/pkg/util"
	"os"
	"strconv"
)

const (
	Port               = 8080
	TecnomagiaSubforum = "f132-tecnomagia"
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
	sessionKey := os.Getenv("DOBBY_SESSION_KEY")

	// Connects to TursoDB and gets config table
	tursoConnectLine := fmt.Sprintf(tursoDbUrl + "?authToken=" + tursoDbToken)
	store := storage.NewStore(tursoConnectLine)

	// Create fables if needed
	configApi := model.NewConfigApi(model.Config{}, *store)
	potionSubApi := model.NewPotionSubApi(model.PotionSub{}, *store)
	potionThrApi := model.NewPotionThreadApi(model.PotionThread{}, *store)
	creationChamberSubApi := model.NewCreationChamberSubApi(model.CreationChamberSub{}, *store)
	userApi := model.NewUserApi(model.User{}, *store)
	announcementApi := model.NewAnnouncementApi(model.Announcement{}, *store)

	err = configApi.CreateInitialConfigTable()
	util.Panic(err)
	err = potionSubApi.CreateInitialPotionSubTable()
	util.Panic(err)
	err = potionThrApi.CreateInitialPotionThreadTable()
	util.Panic(err)
	err = creationChamberSubApi.CreateInitialCreationChamberSubTable()
	util.Panic(err)
	err = userApi.CreateInitialUserTable()
	util.Panic(err)
	err = announcementApi.CreateAnnouncementTable()
	util.Panic(err)

	//Initial user if needed
	userList, err := userApi.GetAllUser()
	util.Panic(err)

	if len(userList) == 0 {
		user := model.User{
			Username:    "Arikel McDowell",
			Active:      true,
			Title:       "Sirena!",
			Permissions: string(model.PermissionAdmin + "," + model.PermissionPotions + "," + model.PermissionCreationChamber),
		}
		fmt.Printf("Inserting test user: \n %s", util.MarshalJsonPretty(user))
		err = userApi.InsertUser(user)
		util.Panic(err)
	}

	// Gets config
	configTable, err := configApi.GetConfig()
	util.Panic(err)
	fmt.Printf("Config: \n %s", util.MarshalJsonPretty(configTable))
	fmt.Println("Server starting at: http://localhost:" + serverPort)

	// Starts the server
	app := echo.New()
	app.Use(slogecho.New(mylogger.GetLogger()))
	app.Static("/assets", "/internal/app/dobby-server/assets")

	// Generate new sessionKey to invalidate all previous sessions
	sessionKey, err = util.GenerateRandomKey(32)
	util.Panic(err)
	app.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionKey))))

	handler.SetupRoutes(app, &configTable, store)

	port := Port
	for {
		util.LongPrintlnPrintln("Server starting at: http://localhost:" + strconv.Itoa(port))
		err := app.Start(":" + strconv.Itoa(port))
		if err != nil {
			if util.IsPortInUse(port) {
				util.LongPrintlnPrintln("Port " + strconv.Itoa(port) + " in use, trying next port...")
				port++
			} else {
				util.LongPrintlnPrintln("Failed to start server: %v", err)
			}
		} else {
			break
		}
	}

	util.Panic(err)
}
