package main

import (
	"fmt"
	"localdev/dobby-server/internal/app/dobby-server/model"
	clientConfig "localdev/dobby-server/internal/app/gobs-client/config"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
)

const (
	pathToConfig = "cmd/gobs-client/conf.json"
)

var ToolArray *[]tool.Tool

func main() {
	// LOAD CONFIGS
	ToolArray = &[]tool.Tool{}
	conf := clientConfig.ReadConfigFile(pathToConfig)
	serverConfig := model.Config{
		BaseUrl:            "https://www.hogwartsrol.com/",
		GSheetTokenFile:    "",
		GSheetCredFile:     "",
		GSheetModeracionId: "",
	}
	fmt.Printf("Config: \n %s", util.MarshalJsonPretty(conf))

	// LOGIN
	for _, user := range conf.Users {
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
		*ToolArray = append(*ToolArray, *o)
	}

	// POST A NEW THREAD FOR EACH USER
	/*
		for _, o := range *ToolArray {
			// POST NEW THREAD
			timeNow := time.Now()
			todayDate := timeNow.Format("02/01/2006")
			subject := "Posteando desde Dobby | " + todayDate
			subforumId := "132" // f132-tecnomagia
			message := `probando!!!  [roll="Ajedrez mágico"][/roll]`
			_, err := o.PostNewThread(subforumId, subject, message, false, false, true)
			util.Panic(err)
		}
	*/
}
