package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"os"
)

const (
	loginUsername = "Desarrollo"
	loginPassword = "programación2055"
	inputCsvFile  = "testUser.csv"
	csvDelimiter  = ','
)

type session struct {
	Tool *tool.Tool
	Conf config
}

type config struct {
	BaseUrl string `json:"baseUrl"`
}

var CurrentUser = ""

type User struct {
	id               string
	userUrl          string
	username         string
	house            string
	lifeStatus       string
	year             string
	side             string
	race             string
	sex              string
	legionMember     string
	age              string
	attack           string
	defense          string
	galeons          string
	title            string
	bloodStatus      string
	inventoryOneHtml string
	inventoryTwoHtml string
}

func main() {
	// FORUM LOGIN
	s := &session{}
	var o *tool.Tool
	serverConfig := model.Config{
		BaseUrl:         "https://www.hogwartsrol.com/",
		GSheetTokenFile: "",
		GSheetCredFile:  "",
	}
	client, loginResponse := tool.LoginAndGetCookies(loginUsername, loginPassword)
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
	s.Conf.BaseUrl = serverConfig.BaseUrl
	s.Tool = o

	// READ CSV FILE

	csvInputFile, err := os.Open(inputCsvFile)
	util.Panic(err)
	reader := csv.NewReader(bufio.NewReader(csvInputFile))
	reader.Comma = csvDelimiter

	isHeaderLine := true
	processedLine := 0

	var users []User
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		util.Panic(err)
		if isHeaderLine {
			isHeaderLine = false
		} else {
			processedLine++
			id := line[0]
			users = append(users, User{
				id:      id,
				userUrl: "https://www.hogwartsrol.com/" + id,
			})

			//GET USER INFO FROM FORUM
			profileHtml := s.Tool.GetUserProfile(id)

			profile := parser.ProfileGetProfile(profileHtml)
			fmt.Println(fmt.Sprintf("%s\n", util.MarshalJsonPretty(profile)))

		}
	}

}
