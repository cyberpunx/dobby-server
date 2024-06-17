package main

import (
	"bufio"
	"flag"
	"fmt"
	"localdev/dobby-server/internal/app/dobby-server/model"
	gobsconfig "localdev/dobby-server/internal/app/gobs-client/config"
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	pathToConfig    = "conf_gobs.json"
	gobsSubject     = "<gobs>"
	gobsFile        = "gobs.txt"
	spamControlWait = 11
)

type session struct {
	Tool *tool.Tool
	Conf gobsconfig.Config
	User gobsconfig.User
}

var SessionArray *[]session

var CurrentUser = ""

func main() {
	name := flag.String("name", "", "Username")
	confPath := flag.String("conf", pathToConfig, "Path to config file")
	flag.Parse()

	if *name == "" {
		fmt.Println("Please provide a username")
		return
	}
	CurrentUser = *name

	// LOAD CONFIGS
	SessionArray = &[]session{}
	conf := gobsconfig.ReadConfigFile(*confPath)
	serverConfig := model.Config{
		BaseUrl:            "https://www.hogwartsrol.com/",
		GSheetTokenFile:    "",
		GSheetCredFile:     "",
		GSheetModeracionId: "",
	}
	fmt.Printf("Config: \n %s", util.MarshalJsonPretty(conf))

	// LOGIN EVERY USER
	for _, user := range conf.Users {
		s := &session{}
		var o *tool.Tool
		client, loginResponse := tool.LoginAndGetCookies(user.Username, user.Password)
		if !*loginResponse.Success {
			fmt.Println("Usuario y/o Contraseña incorrectos")
		} else {
			o = tool.NewTool(&serverConfig, client, nil, nil)
			secret1, secret2, err := o.GetPostSecrets()
			if err != nil {
				fmt.Println("Es posible que el usuario no tenga permisos en el foro / error al obtener secretos")
			}
			o.PostSecret1 = &secret1
			o.PostSecret2 = &secret2
		}
		s.Conf = conf
		s.Tool = o
		s.User = user
		*SessionArray = append(*SessionArray, *s)
	}

	// POST A NEW THREAD FOR EACH USER
	for _, s := range *SessionArray {
		if strings.Contains(strings.ToUpper(s.User.Username), strings.ToUpper(*name)) {
			util.LongPrintlnPrintln("Processing: ", s.User.Username)

			for _, userConfMsg := range s.User.Msg {
				switch userConfMsg.PostDynamic {
				case gobsconfig.GobsPostDynamic:
					err := GobsDynamic(s, userConfMsg)
					util.Panic(err)
				}
			}
		}
	}
}

func GobsDynamic(s session, msgConfig gobsconfig.UserMessage) error {
	//GET SUBFORUM THREADS
	urList := s.Tool.GetThreadsUrlsFromSubforum(msgConfig.SubForumUrl)
	util.LongPrintlnPrintln("Threads: ", urList)
	postedToday := checkIfPostGobTodayWithUser(urList, s.User.Username)
	if postedToday {
		return fmt.Errorf("user has already posted today")
	}

	// POST NEW THREAD
	gobsThread := postGobsNewThread(s, msgConfig)
	gobsThreadUrl := gobsThread.Url

	// GET THREAD AND CHECK IF IT IS A GOBS THREAD
	threadBody := s.Tool.GetThread(gobsThreadUrl)
	thread := s.Tool.ParseThread(threadBody)
	isGobs := parser.PostIsGobstons(thread.Posts[0].Content)

	if isGobs {
		gobsValue := parser.PostGetGobsValue(thread.Posts[0].Content)
		util.LongPrintlnPrintln("GOBS VALUE: ", gobsValue)
		if gobsValue == 0 {
			saveSadFaceToFile(gobsThreadUrl)
		}
		gobsLines := getGobsFileLines()
		if len(gobsLines) == 3 {
			chessMessage := createChessMsg(gobsLines)

			//Sleep X seconds to avoid spamming control
			time.Sleep(spamControlWait * time.Second)

			// POST NEW CHESS THREAD
			chessThread := postChessNewThread(s, msgConfig, chessMessage)
			chessValue := parser.PostGetChessValue(chessThread.Posts[0].Content)
			util.LongPrintlnPrintln("CHESS VALUE: ", chessValue)
			chessLinks := parser.PostGetChessLinks(chessThread.Posts[0].Content)
			util.LongPrintlnPrintln("CHESS LINKS: ", chessLinks)
			eraseGobsFile()
		}
	}

	return nil
}

func saveSadFaceToFile(url string) {
	// Open or create gobs.txt, if exists append at the end
	file, err := os.OpenFile(CurrentUser+"-"+gobsFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	util.Panic(err)
	defer file.Close()

	// Write url
	_, err = file.WriteString(url + "\n")
	util.Panic(err)
}

func getGobsFileLines() []string {
	// Open or create gobs.txt
	file, err := os.OpenFile(CurrentUser+"-"+gobsFile, os.O_RDWR|os.O_CREATE, 0666)
	util.Panic(err)
	defer file.Close()

	// Read file
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func createChessMsg(gobsLines []string) string {
	msg := "[roll=\"Ajedrez mágico\"][/roll]" + "<br><br>"
	if len(gobsLines) == 3 {
		for i, line := range gobsLines {
			msg += "[b]Cara  Triste " + strconv.Itoa(i+1) + ":[/b] [url=https://www.hogwartsrol.com" + line + "]" + strconv.Itoa(i+1) + "º Enlace[/url]" + "<br>"
		}
	}
	return msg
}

func eraseGobsFile() {
	// Open or create gobs.txt
	file, err := os.OpenFile(CurrentUser+"-"+gobsFile, os.O_RDWR|os.O_CREATE, 0666)
	util.Panic(err)
	defer file.Close()

	err = file.Truncate(0)
	util.Panic(err)
}

func postGobsNewThread(s session, MsgConfig gobsconfig.UserMessage) *parser.Thread {
	var subject string

	location, err := time.LoadLocation("America/Mexico_City")
	util.Panic(err)
	currentTime := time.Now().In(location)
	todayDate := currentTime.Format("02/01/2006")
	if MsgConfig.Subject == gobsSubject {
		subject = s.User.Username + " | " + todayDate
	} else {
		subject = MsgConfig.Subject
	}
	subforumId := MsgConfig.SubForumId
	message := MsgConfig.Message
	gobsThread, err := s.Tool.PostNewThread(subforumId, subject, message, false, false, true)
	util.Panic(err)
	util.LongPrintlnPrintln("GOBS THREAD: ", gobsThread.Url)
	return gobsThread

}

func postChessNewThread(s session, MsgConfig gobsconfig.UserMessage, chessMessage string) *parser.Thread {
	var subject string
	location, err := time.LoadLocation("America/Mexico_City")
	util.Panic(err)
	currentTime := time.Now().In(location)
	todayDate := currentTime.Format("02/01/2006")
	if MsgConfig.Subject == gobsSubject {
		subject = s.User.Username + " | " + todayDate
	} else {
		subject = MsgConfig.Subject
	}
	subforumId := "177"
	message := chessMessage
	chessThread, err := s.Tool.PostNewThread(subforumId, subject, message, false, false, true)
	util.Panic(err)
	util.LongPrintlnPrintln("CHESS THREAD: ", chessThread.Url)
	return chessThread
}

func checkIfPostGobTodayWithUser(urlList []string, user string) bool {
	//Checks if user has posted gobs today
	//User Example Hikaru Munetaka
	//URL Example: /t99287-hikaru-munetaka-17-05-2024

	location, err := time.LoadLocation("America/Mexico_City")
	util.Panic(err)
	currentTime := time.Now().In(location)
	datePart := currentTime.Format("02-01-2006")
	userPart := strings.ReplaceAll(strings.ToLower(user), " ", "-")

	for _, url := range urlList {
		//check if url contains user
		if strings.Contains(url, userPart) {
			//check if url contains date
			if strings.Contains(url, datePart) {
				fmt.Println("User has already posted today")
				return true
			}
		}
	}
	return false
}
