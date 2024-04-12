package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/view"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics/potion"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
)

type DobbyHandler struct {
	Tool                  *tool.Tool
	User                  *model.UserSession
	ConfigApi             *model.ConfigApi
	PotionSubApi          *model.PotionSubApi
	PotionThrApi          *model.PotionThreadApi
	CreationChamberSubApi *model.CreationChamberSubApi
}

const (
	loadReportMockup = false
	saveReportMockup = false
)

func (h DobbyHandler) HandleShowLoginForm(c echo.Context) error {
	return render(c, view.Login(*h.User, *h.Tool))
}

func (h DobbyHandler) HandleProcessLoginForm(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	client, loginResponse := tool.LoginAndGetCookies(username, password)
	if !*loginResponse.Success {
		return render(c, view.Login(*h.User, *h.Tool))
	} else {
		h.Tool.Client = client
		secret1, secret2 := h.Tool.GetPostSecrets()
		h.Tool.PostSecret1 = &secret1
		h.Tool.PostSecret2 = &secret2
		userDateFormat := h.Tool.GetUserDateTimeFormat()
		h.User.UserDateFormat = &userDateFormat
		h.User.IsCorrectDateFmt = util.IsUserDateFormatCorrect(userDateFormat, h.Tool.ForumDateTime)
		h.User.Username = &username
		h.User.Initials = loginResponse.Initials
		h.User.LoginDatetime = loginResponse.LoginDatetime
		h.User.IsLoggedIn = true

		if loadReportMockup {
			var report []potion.PotionClubReport
			jsonBytes, err := util.LoadJsonFile("./tmp/potionsReport.json")
			util.Panic(err)
			err = json.Unmarshal(jsonBytes, &report)
			util.Panic(err)

			return render(c, view.Potions(report, *h.User, *h.Tool, "Pociones"))
		}

		return render(c, view.Login(*h.User, *h.Tool))
	}
}

func (h DobbyHandler) HandleLogout(c echo.Context) error {
	h.Tool.Client = nil
	h.User.IsLoggedIn = false
	h.User.Username = nil
	h.User.Initials = nil
	h.User.LoginDatetime = nil
	return render(c, view.Login(*h.User, *h.Tool))
}

func (h DobbyHandler) HandlePotions(c echo.Context) error {
	subForumConfig, err := h.PotionSubApi.GetAllPotionSub()
	util.Panic(err)

	var urls []string
	var timeLimit *int
	var turnLimit *int
	for _, sub := range subForumConfig {
		urls = append(urls, sub.Url)
		timeLimit = &sub.TimeLimit
		turnLimit = &sub.TurnLimit
	}

	potionsReport := h.Tool.ProcessPotionsSubforumList(dynamics.DynamicPotion, &urls, timeLimit, turnLimit)

	if saveReportMockup {
		jsonResponse, err := json.Marshal(potionsReport)
		util.Panic(err)
		//save the json file
		err = util.SaveJsonFile("./tmp/potionsReport.json", jsonResponse)
		util.Panic(err)
	}

	return render(c, view.Potions(potionsReport, *h.User, *h.Tool, "Pociones"))
}

func (h DobbyHandler) HandleCreationChamber(c echo.Context) error {
	subForumConfig, err := h.CreationChamberSubApi.GetAllCreationChamberSub()
	util.Panic(err)

	var urls []string
	var timeLimit *int
	var turnLimit *int
	for _, sub := range subForumConfig {
		urls = append(urls, sub.Url)
		timeLimit = &sub.TimeLimit
		turnLimit = &sub.TurnLimit
	}

	creationChamberReport := h.Tool.ProcessPotionsSubforumList(dynamics.DynamicCreationChamber, &urls, timeLimit, turnLimit)

	if saveReportMockup {
		jsonResponse, err := json.Marshal(creationChamberReport)
		util.Panic(err)
		//save the json file
		err = util.SaveJsonFile("./tmp/creationChamber.json", jsonResponse)
		util.Panic(err)
	}

	return render(c, view.Potions(creationChamberReport, *h.User, *h.Tool, "Cámara de Creación"))
}
