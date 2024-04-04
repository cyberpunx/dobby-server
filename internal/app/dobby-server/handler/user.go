package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/view/user"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"net/http"
)

type UserHandler struct {
	tool *tool.Tool
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h UserHandler) HandleUserShow(c echo.Context) error {
	u := model.User{
		Email:    "aiden@hola.com",
		Username: "Aiden",
		Password: "password123",
	}
	return render(c, user.Show(u))
}

func (h UserHandler) HandleUserLogin(c echo.Context) error {
	r := LoginRequest{}
	err := c.Bind(&r)
	util.Panic(err)

	client, loginResponse := tool.LoginAndGetCookies(r.Username, r.Password)
	if !*loginResponse.Success {
		return c.JSON(http.StatusUnauthorized, loginResponse.Messaage)
	} else {
		h.tool.Client = client
		secret1, secret2 := h.tool.GetPostSecrets()
		h.tool.PostSecret1 = &secret1
		h.tool.PostSecret2 = &secret2

		return render(c, user.Login(secret1, secret2))
	}
}
