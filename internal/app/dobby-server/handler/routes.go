package handler

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo) {
	/*
		group := app.Group("/user")

		group.GET("", h.HandlerShowUsers)
		group.GET("/details/:id", h.HandlerShowUserById)
	*/

	loginHandler := LoginHandler{}
	app.POST("/login", loginHandler.HandleLogin)
}
