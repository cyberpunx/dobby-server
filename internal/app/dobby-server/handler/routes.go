package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/config"
	"localdev/dobby-server/internal/pkg/gsheet"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
)

func SetupRoutes(app *echo.Echo, conf *config.Config) {
	userHandler := UserHandler{
		tool: tool.NewTool(conf, nil, nil),
	}
	sheetService := gsheet.GetSheetService(*userHandler.tool.Config.GSheetTokenFile, *userHandler.tool.Config.GSheetCredFile)
	userHandler.tool.SheetService = sheetService

	app.POST("/login", userHandler.HandleUserLogin)
}
