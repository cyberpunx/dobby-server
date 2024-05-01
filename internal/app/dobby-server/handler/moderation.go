package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/view"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics/potion"
	"localdev/dobby-server/internal/pkg/util"
)

const (
	loadPotionsReportMockup = true
	savePotionsReportMockup = false
)

type ModerationHandler struct {
	h *BaseHandler
}

func (m ModerationHandler) HandlePotions(c echo.Context) error {
	if !m.h.UserSession.HavePermission(model.PermissionPotions) {
		return render(c, view.Home(*m.h.UserSession, *m.h.Tool, "Inicio", "No tienes permisos para ver esta página"))
	}

	subForumConfig, err := m.h.PotionSubApi.GetAllPotionSub()
	util.Panic(err)

	var urls []string
	var timeLimit *int
	var turnLimit *int
	for _, sub := range subForumConfig {
		urls = append(urls, sub.Url)
		timeLimit = &sub.TimeLimit
		turnLimit = &sub.TurnLimit
	}

	potionsReport := m.h.Tool.ProcessPotionsSubforumList(dynamics.DynamicPotion, &urls, timeLimit, turnLimit)

	if savePotionsReportMockup {
		jsonResponse, err := json.Marshal(potionsReport)
		util.Panic(err)
		//save the json file
		err = util.SaveJsonFile("./tmp/potionsReport.json", jsonResponse)
		util.Panic(err)
	}

	return render(c, view.Potions(potionsReport, *m.h.UserSession, *m.h.Tool, potion.PotionNames, "Pociones"))
}

func (a ModerationHandler) HandlePotionNewForm(c echo.Context) error {
	return render(c, view.NewPotionForm(potion.PotionNames))
}

func (m ModerationHandler) HandleNewPotion(c echo.Context) error {
	player1 := c.FormValue("player1")
	player2 := c.FormValue("player2")
	potionName := c.FormValue("potionName")
	mod := *m.h.UserSession.Username
	//subforumUrl := "f98-club-de-pociones" //TODO: Hardcoded for now
	subforumUrl := "f132-tecnomagia" //TODO: Hardcoded for now
	subForumConfig, err := m.h.PotionSubApi.GetAllPotionSub()
	util.Panic(err)
	var turnLimit int
	for _, sub := range subForumConfig {
		turnLimit = sub.TurnLimit
	}
	targetScore := potion.PotionScores[potionName]
	potionMsg := m.h.Tool.GetNewPotionMessage(dynamics.DynamicPotion, player1, player2, potionName, mod, subforumUrl, turnLimit, targetScore)
	return render(c, view.NewPotionMsg(potionMsg))
}

func (m ModerationHandler) HandleCreationChamber(c echo.Context) error {
	if !m.h.UserSession.HavePermission(model.PermissionCreationChamber) {
		return render(c, view.Home(*m.h.UserSession, *m.h.Tool, "Inicio", "No tienes permisos para ver esta página"))
	}

	subForumConfig, err := m.h.CreationChamberSubApi.GetAllCreationChamberSub()
	util.Panic(err)

	var urls []string
	var timeLimit *int
	var turnLimit *int
	for _, sub := range subForumConfig {
		urls = append(urls, sub.Url)
		timeLimit = &sub.TimeLimit
		turnLimit = &sub.TurnLimit
	}

	creationChamberReport := m.h.Tool.ProcessPotionsSubforumList(dynamics.DynamicCreationChamber, &urls, timeLimit, turnLimit)

	if savePotionsReportMockup {
		jsonResponse, err := json.Marshal(creationChamberReport)
		util.Panic(err)
		//save the json file
		err = util.SaveJsonFile("./tmp/creationChamber.json", jsonResponse)
		util.Panic(err)
	}

	return render(c, view.Potions(creationChamberReport, *m.h.UserSession, *m.h.Tool, potion.PotionNames, "Cámara de Creación"))
}
