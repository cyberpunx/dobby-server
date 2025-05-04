package main

import (
	"fmt"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
)

const (
	loginUsername = "Desarrollo"
	loginPassword = "programación2055"
)

func main() {
	// FORUM LOGIN
	s := forumLogin(loginUsername, loginPassword)

}

func forumLogin(username, password string) *session {
	// FORUM LOGIN
	s := &session{}
	var o *tool.Tool
	serverConfig := model.Config{
		BaseUrl:         "https://www.hogwartsrol.com/",
		GSheetTokenFile: "",
		GSheetCredFile:  "",
	}
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
	s.Conf.BaseUrl = serverConfig.BaseUrl
	s.Tool = o

	return s
}
