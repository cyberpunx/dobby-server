package tool

import (
	"google.golang.org/api/sheets/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/storage"
	"localdev/dobby-server/internal/pkg/util"
	"net/http"
	"time"
)

type LoginRequest struct {
	User *string `json:"user"`
	Pass *string `json:"pass"`
}

type LoginResponse struct {
	Success       *bool      `json:"success"`
	Messaage      *string    `json:"message"`
	Username      *string    `json:"username"`
	Initials      *string    `json:"initials"`
	LoginDatetime *time.Time `json:"datetime"`
}

type Tool struct {
	Config        *model.Config
	Client        *http.Client
	SheetService  *sheets.Service
	ForumDateTime time.Time
	Store         *storage.Store
	PostSecret1   *string
	PostSecret2   *string
}

func NewTool(config *model.Config, client *http.Client, gSheetService *sheets.Service, store *storage.Store) *Tool {
	forumDateTime, err := util.GetTimeFromTimeZone("America/Mexico_City")
	util.Panic(err)
	return &Tool{
		Config:        config,
		Client:        client,
		ForumDateTime: forumDateTime,
		SheetService:  gSheetService,
		Store:         store,
	}
}
