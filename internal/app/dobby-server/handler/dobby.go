package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"net/http"
)

type DobbyHandler struct {
	Tool *tool.Tool
	User *model.User
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h DobbyHandler) HandleDobbyLogin(c echo.Context) error {
	r := LoginRequest{}
	err := c.Bind(&r)
	util.Panic(err)

	client, loginResponse := tool.LoginAndGetCookies(r.Username, r.Password)
	if !*loginResponse.Success {
		return c.JSON(http.StatusUnauthorized, loginResponse.Messaage)
	} else {
		h.Tool.Client = client
		secret1, secret2 := h.Tool.GetPostSecrets()
		h.Tool.PostSecret1 = &secret1
		h.Tool.PostSecret2 = &secret2
		h.User = &model.User{
			Username: &r.Username,
			Initials: loginResponse.Initials,
			Datetime: loginResponse.Datetime,
		}
		return c.JSON(http.StatusOK, h.User)
	}
}

func (h DobbyHandler) HandleDobbyPotions(c echo.Context) error {
	subForumConfig := h.Tool.Store.GetPotionSubforum()

	var urls []string
	var timeLimit *int
	var turnLimit *int
	for _, sub := range *subForumConfig {
		urls = append(urls, *sub.Url)
		timeLimit = sub.TimeLimit
		turnLimit = sub.TurnLimit
	}

	h.Tool.ProcessPotionsSubforumList(&urls, timeLimit, turnLimit)

	return c.JSON(http.StatusOK, "Potions processed")
}
