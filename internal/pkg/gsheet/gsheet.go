package gsheet

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"localdev/dobby-server/internal/pkg/util"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

const (
	SheetIdModeration            = "13CCYZ4veljB6ItPNHdvxvClBZJaC1w-QMkq-H5btR74"
	SheetPotionsRangeDaysOff     = "Permisos Pociones!A269:B"
	SheetPotionsRangePlayerBonus = "Logros Pociones!A2:B"

	SheetIdCreationChamber               = "11ts40QG3Bf-909uuzS_ayC5i09YfX5gP_ExWExSiuCY"
	SheetCreationChamberRangeDaysOff     = "Permisos!A50:B"
	SheetCreationChamberRangePlayerBonus = "Permisos!C1:D2"

	SheetIdLogs      = "1JpvrhEFvrasUnL6qDma86uqzyIfAj8Cra5QM60Us4Jo"
	SheetRangeLogins = "Logins!A2:B"
)

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Create a random state token
	state := "st" + string(rand.New(rand.NewSource(time.Now().UnixNano())).Int63())

	// Create a channel to receive the authorization code
	codeChan := make(chan string)

	// Start a web server to listen on the callback URL
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check the state token
		if r.URL.Query().Get("state") != state {
			http.Error(w, "State token does not match", http.StatusBadRequest)
			return
		}
		// Send the code to the channel
		codeChan <- r.URL.Query().Get("code")
		fmt.Fprintf(w, "Authorization complete, you can close this window.")
		server.Shutdown(context.Background())
	})

	// Open the authorization URL in the user's browser
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	util.LongPrintlnPrintln("Opening browser for authorization: %s\n", authURL)
	openbrowser(authURL)

	// Start the server and wait for the auth code
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			util.Panic(err)
		}
	}()
	code := <-codeChan

	// Exchange the code for a token
	tok, err := config.Exchange(context.Background(), code)
	util.Panic(err)
	return tok
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	util.Panic(err)
}
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	util.LongPrintlnPrintln("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	util.Panic(err)
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetClient(tokFile string, ctx context.Context, config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		//get token from environment variable GSHEET_TOKEN
		tokEnvValue := os.Getenv("GSHEET_TOKEN")
		if len(tokEnvValue) > 0 {
			tok = &oauth2.Token{}
			err = json.Unmarshal([]byte(tokEnvValue), tok)
			if err != nil {
				tok = getTokenFromWeb(config)
				saveToken(tokFile, tok)
			}
		}
	}

	return config.Client(ctx, tok)
}

func ReadCredentials(credPath string) ([]byte, error) {
	b, err := ioutil.ReadFile(credPath)
	if err != nil {
		fmt.Println("unable to read client secret file: %w")
		fmt.Println("attempt to read client secret from environment variable")
		b = []byte(os.Getenv("GSHEET_CLIENT_SECRET"))
		if len(b) == 0 {
			return nil, fmt.Errorf("unable to read client secret from environment variable")
		}
	}
	return b, nil
}

func ReadSheetData(srv *sheets.Service, spreadsheetId, readRange string) ([][]interface{}, error) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}
	return resp.Values, nil
}

func WriteSheetData(srv *sheets.Service, spreadsheetId, writeRange string, rowData []interface{}) error {
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{rowData},
	}

	_, err := srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to write data to sheet: %w", err)
	}
	return nil
}

func FindNextAvailableRow(srv *sheets.Service, spreadsheetId, readRange string) (int, error) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return 0, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	// Extract the start row index from the readRange.
	var startIndex int
	n, err := fmt.Sscanf(readRange, "%*[A-Za-z]!A%d", &startIndex)
	if err != nil || n != 1 {
		startIndex = 2 // Default to 2 if unable to parse range or range does not include a start index.
	}

	// Iterate over the rows to find the first blank row.
	for i, row := range resp.Values {
		if isRowEmpty(row) {
			// Return the index of the first empty row.
			return startIndex + i, nil
		}
	}

	// If no empty row is found, return the row number after the last row.
	return startIndex + len(resp.Values), nil
}

// Helper function to check if a row is empty.
func isRowEmpty(row []interface{}) bool {
	for _, cell := range row {
		if cell != nil && cell != "" {
			return false
		}
	}
	return true
}

func DisplayData(data [][]interface{}) {
	if len(data) == 0 {
		util.LongPrintlnPrintln("No data found.")
	} else {
		util.LongPrintlnPrintln("Data:")
		for _, row := range data {
			util.LongPrintlnPrintln("%s\n", row)
		}
	}
}

func GetSheetService(tokFile, credPath string) *sheets.Service {
	ctx := context.Background()
	// Read Credentials
	credentials, err := ReadCredentials(credPath)
	util.Panic(err)
	//credentials := []byte(ClientSecret)

	// Configure OAuth2 Client
	gconfig, err := google.ConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	util.Panic(err)
	client := GetClient(tokFile, ctx, gconfig)

	// Create Sheets Service
	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	util.Panic(err)
	return service
}

func ParseDayOff(rows [][]interface{}) []DayOff {
	var daysOff []DayOff
	for _, row := range rows {
		//skip first row because it's the header
		//if i == 0 {
		//	continue
		//}
		if len(row) == 2 {
			date, err := time.Parse("2/01/2006", row[1].(string)) //format DD/MM/YYYY
			util.Panic(err)
			dayOff := DayOff{
				Username: row[0].(string),
				Date:     date,
			}
			daysOff = append(daysOff, dayOff)
		}
	}
	return daysOff
}

func ParsePlayerBonus(rows [][]interface{}) []PlayerBonus {
	var playerBonus []PlayerBonus
	for _, row := range rows {
		//skip first row because it's the header
		//if i == 0 {
		//	continue
		//}
		if len(row) == 2 {
			bonusValue, err := strconv.Atoi(row[1].(string))
			util.Panic(err)
			pBonus := PlayerBonus{
				Username: row[0].(string),
				Bonus:    bonusValue,
			}
			playerBonus = append(playerBonus, pBonus)
		}
	}
	return playerBonus
}

func FindDayOffForUser(daysOff *[]DayOff, username string) *DayOff {
	for _, dayOff := range *daysOff {
		dayOffUsername := util.TrimAndToLower(dayOff.Username)
		username = util.TrimAndToLower(username)
		if dayOffUsername == username {
			return &dayOff
		}
	}
	return nil
}

func FindDayOffForPLayerBetweenDates(daysOff *[]DayOff, username string, startDate, endDate time.Time) *DayOff {
	for _, dayOff := range *daysOff {
		dayOffUsername := util.TrimAndToLower(dayOff.Username)
		username = util.TrimAndToLower(username)
		if dayOffUsername == username && dayOff.Date.After(startDate) && dayOff.Date.Before(endDate) {
			return &dayOff
		}
	}
	return nil
}

func GetPlayerBonusForUsername(playerBonus *[]PlayerBonus, username string) int {
	for _, pb := range *playerBonus {
		pbUser := util.TrimAndToLower(pb.Username)
		username = util.TrimAndToLower(username)
		if pbUser == username {
			return pb.Bonus
		}
	}
	return 0
}
