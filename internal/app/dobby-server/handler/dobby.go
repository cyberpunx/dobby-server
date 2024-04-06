package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/view"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
)

type DobbyHandler struct {
	Tool *tool.Tool
	User *model.User
}

func (h DobbyHandler) HandleShowLoginForm(c echo.Context) error {
	return render(c, view.Login(h.User))
}

func (h DobbyHandler) HandleProcessLoginForm(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	client, loginResponse := tool.LoginAndGetCookies(username, password)
	if !*loginResponse.Success {
		return render(c, view.Login(h.User))
	} else {
		h.Tool.Client = client
		secret1, secret2 := h.Tool.GetPostSecrets()
		h.Tool.PostSecret1 = &secret1
		h.Tool.PostSecret2 = &secret2
		h.User.Username = &username
		h.User.Initials = loginResponse.Initials
		h.User.Datetime = loginResponse.Datetime
		h.User.IsLoggedIn = true
		return render(c, view.Login(h.User))
	}
}

func (h DobbyHandler) HandleLogout(c echo.Context) error {
	h.Tool.Client = nil
	h.User.IsLoggedIn = false
	h.User.Username = nil
	h.User.Initials = nil
	h.User.Datetime = nil
	return render(c, view.Login(h.User))
}

func (h DobbyHandler) HandlePotions(c echo.Context) error {
	subForumConfig := h.Tool.Store.GetPotionSubforum()

	var urls []string
	var timeLimit *int
	var turnLimit *int
	for _, sub := range *subForumConfig {
		urls = append(urls, *sub.Url)
		timeLimit = sub.TimeLimit
		turnLimit = sub.TurnLimit
	}

	potionsReport := h.Tool.ProcessPotionsSubforumList(&urls, timeLimit, turnLimit)

	return render(c, view.Potions(potionsReport))
}
