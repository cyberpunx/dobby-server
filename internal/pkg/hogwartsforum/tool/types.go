package tool

import (
	"google.golang.org/api/sheets/v4"
	"localdev/dobby-server/internal/app/dobby-server/config"
	"localdev/dobby-server/internal/pkg/util"
	"net/http"
	"time"
)

type LoginRequest struct {
	User *string `json:"user"`
	Pass *string `json:"pass"`
}

type LoginResponse struct {
	Success  *bool      `json:"success"`
	Messaage *string    `json:"message"`
	Username *string    `json:"username"`
	Initials *string    `json:"initials"`
	Datetime *time.Time `json:"datetime"`
}

type Tool struct {
	Config        *config.Config
	Client        *http.Client
	SheetService  *sheets.Service
	ForumDateTime time.Time
	PostSecret1   *string
	PostSecret2   *string
}

func NewTool(config *config.Config, client *http.Client, gSheetService *sheets.Service) *Tool {
	forumDateTime, err := util.GetTimeFromTimeZone("America/Mexico_City")
	util.Panic(err)
	return &Tool{
		Config:        config,
		Client:        client,
		ForumDateTime: forumDateTime,
		SheetService:  gSheetService,
	}
}
