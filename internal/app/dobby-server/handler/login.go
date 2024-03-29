package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/view/login"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
)

type LoginHandler struct{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h LoginHandler) HandleLogin(c echo.Context) error {
	r := LoginRequest{}
	err := c.Bind(&r)

	tool.LoginAndGetCookies(r.Username, r.Password)

	util.Panic(err)
	return render(c, login.Login())
}
