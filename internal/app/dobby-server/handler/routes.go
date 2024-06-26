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
	AnnouncementApi       *model.AnnouncementApi
}

func SetupRoutes(app *echo.Echo, conf *model.Config, store *storage.Store) {
	handler := BaseHandler{
		Tool: tool.NewTool(conf, nil, gsheet.GetSheetService(conf.GSheetTokenFile, conf.GSheetCredFile), store),
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
		AnnouncementApi:       model.NewAnnouncementApi(model.Announcement{}, *store),
	}

	modHandler := ModerationHandler{&handler}
	moderationGroup := app.Group("/moderation")
	moderationGroup.GET("/potions", modHandler.HandlePotions)
	moderationGroup.GET("/potion/new", modHandler.HandlePotionNewForm)
	moderationGroup.POST("/potion/new", modHandler.HandleNewPotion)
	moderationGroup.GET("/creationchamber", modHandler.HandleCreationChamber)

	loginHandler := LoginHandler{&handler}
	app.POST("/login", loginHandler.HandleProcessLoginForm)
	app.GET("/logout", loginHandler.HandleLogout)

	pagesHandler := PagesHandler{&handler}
	app.GET("/", pagesHandler.HandleHome)
	app.GET("/timecheck", pagesHandler.HandleTimeCheckForm)
	app.POST("/timecheck", pagesHandler.HandleTimeCheck)

	adminHandler := AdminHandler{&handler}
	adminGroup := app.Group("/admin")
	adminGroup.GET("/user/list", adminHandler.HandleUserList)
	adminGroup.GET("/user/:id/edit", adminHandler.HandleUserEdit)
	adminGroup.DELETE("/user/:id", adminHandler.HandleUserDelete)
	adminGroup.PUT("/user/:id", adminHandler.HandleUserUpdate)
	adminGroup.GET("/user/:id", adminHandler.HandleUserView)
	adminGroup.GET("/user/new", adminHandler.HandleUserNewForm)
	adminGroup.POST("/user/new", adminHandler.HandleUserNew)

	adminGroup.GET("/announcement/list", adminHandler.HandleAnnouncementList)
	adminGroup.GET("/announcement/:id/edit", adminHandler.HandleAnnouncementEdit)
	adminGroup.DELETE("/announcement/:id", adminHandler.HandleAnnouncementDelete)
	adminGroup.PUT("/announcement/:id", adminHandler.HandleAnnouncementUpdate)
	adminGroup.GET("/announcement/:id", adminHandler.HandleAnnouncementView)
	adminGroup.GET("/announcement/new", adminHandler.HandleAnnouncementNewForm)
	adminGroup.POST("/announcement/new", adminHandler.HandleAnnouncementNew)

}
