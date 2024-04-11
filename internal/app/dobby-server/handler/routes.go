package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/storage"
	"localdev/dobby-server/internal/pkg/gsheet"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
)

func SetupRoutes(app *echo.Echo, conf *model.Config, store *storage.Store) {
	dobbyGroup := app.Group("/dobby")
	dobbyHandler := DobbyHandler{
		Tool: tool.NewTool(conf, nil, nil, store),
		User: &model.User{
			IsLoggedIn:    false,
			Username:      nil,
			Initials:      nil,
			LoginDatetime: nil,
		},
		ConfigApi:             model.NewConfigApi(model.Config{}, *store),
		PotionSubApi:          model.NewPotionSubApi(model.PotionSub{}, *store),
		PotionThrApi:          model.NewPotionThreadApi(model.PotionThread{}, *store),
		CreationChamberSubApi: model.NewCreationChamberSubApi(model.CreationChamberSub{}, *store),
	}
	sheetService := gsheet.GetSheetService(dobbyHandler.Tool.Config.GSheetTokenFile, dobbyHandler.Tool.Config.GSheetCredFile)
	dobbyHandler.Tool.SheetService = sheetService
	dobbyGroup.POST("/login", dobbyHandler.HandleProcessLoginForm)
	dobbyGroup.GET("/logout", dobbyHandler.HandleLogout)
	dobbyGroup.GET("/potions", dobbyHandler.HandlePotions)
	dobbyGroup.GET("/creationchamber", dobbyHandler.HandleCreationChamber)
	app.GET("/", dobbyHandler.HandleShowLoginForm)

}
