package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/view"
	"localdev/dobby-server/internal/pkg/util"
)

type PagesHandler struct {
	h *BaseHandler
}

func (p PagesHandler) HandleHome(c echo.Context) error {
	p.h.UserSession, p.h.Tool = GetCurrentSessionAndTool(c, p.h.Tool)

	announcementList, err := p.h.AnnouncementApi.GetAllAnnouncement()
	util.Panic(err)

	if p.h.UserSession.IsLoggedIn {
		return render(c, view.Home(*p.h.UserSession, *p.h.Tool, "Inicio", "", &announcementList))
	} else {
		return render(c, view.Login("", *p.h.UserSession, *p.h.Tool))
	}
}
