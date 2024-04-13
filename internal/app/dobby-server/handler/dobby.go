package handler

import (
	"encoding/json"
	"fmt"
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
	UserSession           *model.UserSession
	ConfigApi             *model.ConfigApi
	PotionSubApi          *model.PotionSubApi
	PotionThrApi          *model.PotionThreadApi
	CreationChamberSubApi *model.CreationChamberSubApi
	UserApi               *model.UserApi
}

const (
	loadReportMockup = false
	saveReportMockup = false
)

func (h DobbyHandler) HandleHome(c echo.Context) error {
	if h.UserSession.IsLoggedIn {
		return render(c, view.Home(*h.UserSession, *h.Tool, "Inicio", ""))
	} else {
		return render(c, view.Login("", *h.UserSession, *h.Tool))
	}
}

func (h DobbyHandler) HandleShowLoginForm(c echo.Context) error {
	return render(c, view.Login("", *h.UserSession, *h.Tool))
}

func (h DobbyHandler) HandleProcessLoginForm(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	client, loginResponse := tool.LoginAndGetCookies(username, password)
	if !*loginResponse.Success {
		return render(c, view.Login("Usuario y/o Contraseña incorrectos", *h.UserSession, *h.Tool))
	} else {
		// UserSession is logged in at Forum
		h.Tool.Client = client
		secret1, secret2, err := h.Tool.GetPostSecrets()
		if err != nil {
			return render(c, view.Login("Es posible que el usuario no tenga permisos en el foro / error al obtener secretos", *h.UserSession, *h.Tool))
		}
		h.Tool.PostSecret1 = &secret1
		h.Tool.PostSecret2 = &secret2

		//Check if the user has permissions on Dobby
		user, err := h.UserApi.GetUserByUsername(username)
		if err != nil || user == nil {
			return render(c, view.Login("No tienes permisos para Dobby", *h.UserSession, *h.Tool))
		}

		h.UserSession.User = user
		h.UserSession.Permissions = user.GetUserPermissions()
		userDateFormat := h.Tool.GetUserDateTimeFormat()
		h.UserSession.UserDateFormat = &userDateFormat
		h.UserSession.IsCorrectDateFmt = util.IsUserDateFormatCorrect(userDateFormat, h.Tool.ForumDateTime)
		h.UserSession.Username = &username
		h.UserSession.Initials = loginResponse.Initials
		h.UserSession.LoginDatetime = loginResponse.LoginDatetime
		h.UserSession.IsLoggedIn = true
		fmt.Printf("UserSession: \n %s", util.MarshalJsonPretty(h.UserSession))

		if loadReportMockup {
			var report []potion.PotionClubReport
			jsonBytes, err := util.LoadJsonFile("./tmp/potionsReport.json")
			util.Panic(err)
			err = json.Unmarshal(jsonBytes, &report)
			util.Panic(err)

			return render(c, view.Potions(report, *h.UserSession, *h.Tool, "Pociones"))
		}

		return render(c, view.Home(*h.UserSession, *h.Tool, "Inicio", ""))
	}
}

func (h DobbyHandler) HandleLogout(c echo.Context) error {
	h.Tool.Client = nil
	h.UserSession.IsLoggedIn = false
	h.UserSession.Username = nil
	h.UserSession.Initials = nil
	h.UserSession.LoginDatetime = nil
	h.UserSession.Permissions = nil
	h.UserSession.User = nil
	h.UserSession.UserDateFormat = nil
	h.UserSession.IsCorrectDateFmt = false
	return render(c, view.Login("", *h.UserSession, *h.Tool))
}

func (h DobbyHandler) HandlePotions(c echo.Context) error {
	if !h.UserSession.HavePermission(model.PermissionPotions) {
		return render(c, view.Home(*h.UserSession, *h.Tool, "Inicio", "No tienes permisos para ver esta página"))
	}

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

	return render(c, view.Potions(potionsReport, *h.UserSession, *h.Tool, "Pociones"))
}

func (h DobbyHandler) HandleCreationChamber(c echo.Context) error {
	if !h.UserSession.HavePermission(model.PermissionCreationChamber) {
		return render(c, view.Home(*h.UserSession, *h.Tool, "Inicio", "No tienes permisos para ver esta página"))
	}

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

	return render(c, view.Potions(creationChamberReport, *h.UserSession, *h.Tool, "Cámara de Creación"))
}
