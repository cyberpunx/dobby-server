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
	pathToConfig = "conf.json"
	gobsSubject  = "<gobs>"
	gobsFile     = "gobs.txt"
)

type session struct {
	Tool *tool.Tool
	Conf gobsconfig.Config
	User gobsconfig.User
}

var SessionArray *[]session

func main() {
	name := flag.String("name", "", "Username")
	confPath := flag.String("conf", pathToConfig, "Path to config file")
	flag.Parse()

	if *name == "" {
		fmt.Println("Please provide a username")
		return
	}

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

	// LOGIN
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
			fmt.Println("Processing: ", s.User.Username)
			var subject string

			// POST NEW THREAD
			timeNow := time.Now()
			todayDate := timeNow.Format("02/01/2006")
			if s.User.Msg.Subject == gobsSubject {
				subject = s.User.Username + " | " + todayDate
			} else {
				subject = s.User.Msg.Subject
			}
			subforumId := s.User.Msg.SubForumId
			message := s.User.Msg.Message
			gobsThread, err := s.Tool.PostNewThread(subforumId, subject, message, false, false, true)
			util.Panic(err)
			fmt.Println("GOBS THREAD: ", gobsThread.Url)
			gobsThreadUrl := gobsThread.Url

			// GET THREAD
			testSession := (*SessionArray)[0]
			threadBody := testSession.Tool.GetThread(gobsThreadUrl)
			thread := testSession.Tool.ParseThread(threadBody)
			isGobs := parser.PostIsGobstons(thread.Posts[0].Content)
			if isGobs {
				gobsValue := parser.PostGetGobsValue(thread.Posts[0].Content)
				if gobsValue == 0 {
					saveSadFaceToFile(gobsThreadUrl)
				}
				gobsLines := getGobsFileLines()
				if len(gobsLines) == 3 {
					msg := createChessMsg(gobsLines)

					// POST NEW THREAD
					if s.User.Msg.Subject == gobsSubject {
						subject = s.User.Username + " | " + todayDate
					} else {
						subject = s.User.Msg.Subject
					}
					subforumId = "177"
					message = msg
					chessThread, err := s.Tool.PostNewThread(subforumId, subject, message, false, false, true)
					util.Panic(err)
					fmt.Println("CHESS THREAD: ", chessThread.Url)
					eraseGobsFile()
				}

			}
		}
	}
}

func saveSadFaceToFile(url string) {
	// Open or create gobs.txt, if exists append at the end
	file, err := os.OpenFile(gobsFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	util.Panic(err)
	defer file.Close()

	// Write url
	_, err = file.WriteString(url + "\n")
	util.Panic(err)
}

func getGobsFileLines() []string {
	// Open or create gobs.txt
	file, err := os.OpenFile(gobsFile, os.O_RDWR|os.O_CREATE, 0666)
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
	file, err := os.OpenFile(gobsFile, os.O_RDWR|os.O_CREATE, 0666)
	util.Panic(err)
	defer file.Close()

	err = file.Truncate(0)
	util.Panic(err)
}
