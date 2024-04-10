package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/config"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/storage/tursodb"
	"localdev/dobby-server/internal/pkg/gsheet"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
)

func SetupRoutes(app *echo.Echo, conf *config.Config, store *tursodb.Store) {
	dobbyGroup := app.Group("/dobby")
	dobbyHandler := DobbyHandler{
		Tool: tool.NewTool(conf, nil, nil, store),
		User: &model.User{
			IsLoggedIn: false,
			Username:   nil,
			Initials:   nil,
			Datetime:   nil,
		},
	}
	sheetService := gsheet.GetSheetService(*dobbyHandler.Tool.Config.GSheetTokenFile, *dobbyHandler.Tool.Config.GSheetCredFile)
	dobbyHandler.Tool.SheetService = sheetService
	dobbyGroup.POST("/login", dobbyHandler.HandleProcessLoginForm)
	dobbyGroup.GET("/logout", dobbyHandler.HandleLogout)
	dobbyGroup.GET("/potions", dobbyHandler.HandlePotions)
	dobbyGroup.GET("/creationchamber", dobbyHandler.HandleCreationChamber)
	app.GET("/", dobbyHandler.HandleShowLoginForm)

}
