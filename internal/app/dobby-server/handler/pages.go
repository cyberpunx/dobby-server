package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/view"
)

type PagesHandler struct {
	h *BaseHandler
}

func (p PagesHandler) HandleHome(c echo.Context) error {
	if p.h.UserSession.IsLoggedIn {
		return render(c, view.Home(*p.h.UserSession, *p.h.Tool, "Inicio", ""))
	} else {
		return render(c, view.Login("", *p.h.UserSession, *p.h.Tool))
	}
}
