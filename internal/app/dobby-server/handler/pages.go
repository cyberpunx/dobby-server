package handler

import (
	"fmt"
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

func (p PagesHandler) HandleTimeCheckForm(c echo.Context) error {
	p.h.UserSession, p.h.Tool = GetCurrentSessionAndTool(c, p.h.Tool)
	return render(c, view.TimeCheckForm(*p.h.UserSession, *p.h.Tool))
}

func (p PagesHandler) HandleTimeCheck(c echo.Context) error {
	p.h.UserSession, p.h.Tool = GetCurrentSessionAndTool(c, p.h.Tool)
	threadUrl := c.FormValue("threadUrl")
	fmt.Println("Thread URL: ", threadUrl)
	//Sanitize the threadUrl

	elapsedTime := p.h.Tool.CheckThreadElapsedTime(threadUrl)
	timecheckMsg := "Tiempo desde el último post: " + util.GetElapsedTime(elapsedTime)
	return render(c, view.TimeCheckMsg(timecheckMsg))
}
