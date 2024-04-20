package handler

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/view"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics/potion"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"time"
)

const (
	ByPassForumLogin = false
)

type LoginHandler struct {
	h *BaseHandler
}

func (l LoginHandler) HandleProcessLoginForm(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if ByPassForumLogin {
		user, err := l.h.UserApi.GetUserByUsername(username)
		if err != nil || user == nil {
			return render(c, view.Login("No tienes permisos para Dobby", *l.h.UserSession, *l.h.Tool))
		}
		l.h.UserSession.User = user
		l.h.UserSession.Permissions = user.GetUserPermissions()
		l.h.UserSession.UserDateFormat = util.PStr("19/4/2024, 06:51")
		l.h.UserSession.IsCorrectDateFmt = true
		l.h.UserSession.Username = &username
		l.h.UserSession.Initials = nil
		l.h.UserSession.LoginDatetime = nil
		l.h.UserSession.IsLoggedIn = true
		fmt.Printf("UserSession: \n %s", util.MarshalJsonPretty(l.h.UserSession))

		if loadPotionsReportMockup {
			var report []potion.PotionClubReport
			jsonBytes, err := util.LoadJsonFile("./tmp/potionsReport.json")
			util.Panic(err)
			err = json.Unmarshal(jsonBytes, &report)
			util.Panic(err)

			return render(c, view.Potions(report, *l.h.UserSession, *l.h.Tool, "Pociones"))
		}

		return render(c, view.Home(*l.h.UserSession, *l.h.Tool, "Inicio", ""))
	}

	// First check if user is a Dobby user
	isDobbyUser, user := l.IsDobbyUser(username)
	if !isDobbyUser {
		return render(c, view.Login("No tienes permisos para Dobby", *l.h.UserSession, *l.h.Tool))

	}

	client, loginResponse := tool.LoginAndGetCookies(username, password)
	if !*loginResponse.Success {
		return render(c, view.Login("Usuario y/o Contraseña incorrectos", *l.h.UserSession, *l.h.Tool))
	}

	l.h.Tool.Client = client
	err := l.SetPostSecrets()
	if err != nil {
		return render(c, view.Login("Es posible que el usuario no tenga permisos en el foro / error al obtener secretos", *l.h.UserSession, *l.h.Tool))
	}

	l.SetUserSession(user, loginResponse.Initials, loginResponse.LoginDatetime)
	fmt.Printf("UserSession: \n %s", util.MarshalJsonPretty(l.h.UserSession))

	if loadPotionsReportMockup {
		var report []potion.PotionClubReport
		jsonBytes, err := util.LoadJsonFile("./tmp/potionsReport.json")
		util.Panic(err)
		err = json.Unmarshal(jsonBytes, &report)
		util.Panic(err)

		return render(c, view.Potions(report, *l.h.UserSession, *l.h.Tool, "Pociones"))
	}

	return render(c, view.Home(*l.h.UserSession, *l.h.Tool, "Inicio", ""))
}

func (l LoginHandler) HandleLogout(c echo.Context) error {
	l.h.Tool.Client = nil
	l.h.Tool.PostSecret1 = nil
	l.h.Tool.PostSecret2 = nil
	l.h.UserSession = &model.UserSession{}
	return render(c, view.Login("", *l.h.UserSession, *l.h.Tool))
}

func (l LoginHandler) IsDobbyUser(username string) (bool, *model.User) {
	user, err := l.h.UserApi.GetUserByUsername(username)
	if err != nil || user == nil {
		return false, nil
	}
	return true, user
}

func (l LoginHandler) SetPostSecrets() error {
	secret1, secret2, err := l.h.Tool.GetPostSecrets()
	if err != nil {
		return err
	}
	l.h.Tool.PostSecret1 = &secret1
	l.h.Tool.PostSecret2 = &secret2
	return nil
}

func (l LoginHandler) SetUserSession(user *model.User, initials *string, loginDateTime *time.Time) {
	l.h.UserSession.User = user
	l.h.UserSession.Permissions = user.GetUserPermissions()
	userDateFormat := l.h.Tool.GetUserDateTimeFormat()
	l.h.UserSession.UserDateFormat = &userDateFormat
	l.h.UserSession.IsCorrectDateFmt = util.IsUserDateFormatCorrect(userDateFormat, l.h.Tool.ForumDateTime)
	l.h.UserSession.Username = &user.Username
	l.h.UserSession.Initials = initials
	l.h.UserSession.LoginDatetime = loginDateTime
	l.h.UserSession.IsLoggedIn = true
}
