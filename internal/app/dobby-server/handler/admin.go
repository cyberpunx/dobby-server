package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/view/crud"
	"localdev/dobby-server/internal/pkg/util"
)

type AdminHandler struct {
	h *BaseHandler
}

func (a AdminHandler) HandleUserList(c echo.Context) error {

	userList, err := a.h.UserApi.GetAllUser()
	util.Panic(err)

	return render(c, view.UserList(*a.h.UserSession, *a.h.Tool, userList))
}
