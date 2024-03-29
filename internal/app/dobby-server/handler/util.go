package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, compontent templ.Component) error {
	return compontent.Render(c.Request().Context(), c.Response())
}
