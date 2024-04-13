package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/storage"
	"localdev/dobby-server/internal/pkg/gsheet"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
)

type BaseHandler struct {
	Tool                  *tool.Tool
	UserSession           *model.UserSession
	ConfigApi             *model.ConfigApi
	PotionSubApi          *model.PotionSubApi
	PotionThrApi          *model.PotionThreadApi
	CreationChamberSubApi *model.CreationChamberSubApi
	UserApi               *model.UserApi
}

func SetupRoutes(app *echo.Echo, conf *model.Config, store *storage.Store) {
	handler := BaseHandler{
		Tool: tool.NewTool(conf, nil, nil, store),
		UserSession: &model.UserSession{
			IsLoggedIn:    false,
			Username:      nil,
			Initials:      nil,
			LoginDatetime: nil,
		},
		UserApi:               model.NewUserApi(model.User{}, *store),
		ConfigApi:             model.NewConfigApi(model.Config{}, *store),
		PotionSubApi:          model.NewPotionSubApi(model.PotionSub{}, *store),
		PotionThrApi:          model.NewPotionThreadApi(model.PotionThread{}, *store),
		CreationChamberSubApi: model.NewCreationChamberSubApi(model.CreationChamberSub{}, *store),
	}
	sheetService := gsheet.GetSheetService(handler.Tool.Config.GSheetTokenFile, handler.Tool.Config.GSheetCredFile)
	handler.Tool.SheetService = sheetService

	modHandler := ModerationHandler{&handler}
	moderationGroup := app.Group("/moderation")
	moderationGroup.GET("/potions", modHandler.HandlePotions)
	moderationGroup.GET("/creationchamber", modHandler.HandleCreationChamber)

	loginHandler := LoginHandler{&handler}
	app.POST("/login", loginHandler.HandleProcessLoginForm)
	app.GET("/logout", loginHandler.HandleLogout)

	pagesHandler := PagesHandler{&handler}
	app.GET("/", pagesHandler.HandleHome)
}
