package main

import (
	"fmt"
	"localdev/dobby-server/internal/app/dobby-server/model"
	clientConfig "localdev/dobby-server/internal/app/gobs-client/config"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"time"
)

const (
	username = ""
	password = ""
)

var o *tool.Tool

func main() {
	// LOAD CONFIGS
	conf := clientConfig.NewConfig()
	serverConfig := model.Config{
		BaseUrl:            "https://www.hogwartsrol.com/",
		GSheetTokenFile:    "",
		GSheetCredFile:     "",
		GSheetModeracionId: "",
	}
	fmt.Printf("Config: \n %s", util.MarshalJsonPretty(conf))

	// LOGIN
	client, loginResponse := tool.LoginAndGetCookies(username, password)
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

	// POST NEW THREAD
	timeNow := time.Now()
	todayDate := timeNow.Format("02/01/2006")
	subject := "Arikel McDowell | " + todayDate
	subforumId := "44" // f44-ocio -> id = 44
	message := `probando`
	thread, err := o.PostNewThread(subforumId, subject, message, false, false)
	util.Panic(err)
	fmt.Printf("Thread: \n %s", util.MarshalJsonPretty(thread))
}
