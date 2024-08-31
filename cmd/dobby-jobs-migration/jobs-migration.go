package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"os"
	"strings"
)

const (
	loginUsername              = "Desarrollo"
	loginPassword              = "programación2055"
	csvDelimiter               = ','
	MemberTableOnNewForum      = "smf_members.csv"
	UsersToMigrateFromOldForum = "Datos de usuarios - Migrados HR.csv"
	//UsersToMigrateFromOldForum = "TEST.csv"
	JobsPost = "jobs.html"
)

func main() {

	// FORUM LOGIN
	s := forumLogin(loginUsername, loginPassword)
	oldUsers := LoadOldForumUsersFromCsv(UsersToMigrateFromOldForum, s)
	newUsers := LoadMembersFromCsv(MemberTableOnNewForum)

	//PAIR USERS BASED ON USERNAME
	migratedUsers := PairUsersBasedOnCsvId(oldUsers, newUsers)

	//get JobsPost file contents as string
	htmlStr, err := LeerArchivo(JobsPost)

	newHtml, err := ReplaceLinks(htmlStr, migratedUsers)
	util.Panic(err)
	fmt.Println(newHtml)
}

func LeerArchivo(filename string) (string, error) {
	// Lee el contenido del archivo
	contenido, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Convierte el contenido a string y lo retorna
	return string(contenido), nil
}

func ReplaceLinks(htmlStr string, migratedUsers []MigratedUser) (string, error) {
	// Crea un documento goquery a partir del string HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	newProfileUrl := "https://harrypotterhead.com/foro/index.php?action=profile;u={new_user_id}"

	// Selecciona todas las etiquetas <a>
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		text := item.Text()
		// Busca un parámetro que coincida con el texto de la etiqueta <a>
		for _, user := range migratedUsers {
			if text == user.OldForumUser.Username {
				newUrl := strings.ReplaceAll(newProfileUrl, "{new_user_id}", user.NewForumUser.Id)
				item.SetAttr("href", newUrl)
				item.SetText(user.NewForumUser.Username)
				break
			}
		}
	})

	// Convierte el documento de nuevo a un string HTML
	var result strings.Builder
	if err := goquery.Render(&result, doc.Selection); err != nil {
		return "", err
	}

	return result.String(), nil
}

func replaceProfileUrl(migratedUsers []MigratedUser) {
	file, err := os.Open(JobsPost)
	util.Panic(err)
	defer file.Close()

	//newProfileUrl := "https://harrypotterhead.com/foro/index.php?action=profile;u={new_user_id}"
	//oldProfileUrl := "https://www.hogwartsrol.com/{old_user_id}"

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

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

func LoadOldForumUsersFromCsv(inputFile string, sessionLoggedIn *session) []OldForumUser {
	// READ USERS CSV FILE
	csvInputFile, err := os.Open(inputFile)
	util.Panic(err)
	reader := csv.NewReader(bufio.NewReader(csvInputFile))
	reader.Comma = csvDelimiter

	isHeaderLine := true
	processedLine := 0

	var oldForumUsers []OldForumUser
	//NEW MEMBERS
	newMembers := LoadMembersFromCsv(MemberTableOnNewForum)

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
			idNewForum := line[1]
			//GET USER INFO FROM FORUM
			profileHtml := sessionLoggedIn.Tool.GetUserProfile(id)

			profile := parser.ProfileGetProfile(profileHtml)
			findMember := SearchMemberById(newMembers, idNewForum)
			if profile.Username != findMember.Username {
				fmt.Println(profile.Username + " - " + findMember.Username)
			}

			totalItems := profile.Inventory.Items

			for _, item := range totalItems {
				if strings.Contains(item.Name, " ") {
					fmt.Println("ParsedItem con espacio en blanco: ", item.Name)
				}
			}

			//CREATE OLD FORUM USER
			oldForumUsers = append(oldForumUsers, OldForumUser{
				Id:         id,
				HPHId:      idNewForum,
				Username:   profile.Username,
				Profile:    profile,
				TotalItems: totalItems,
			})
		}
	}
	return oldForumUsers
}

func SearchMemberById(members []NewForumUser, id string) *NewForumUser {
	for _, member := range members {
		if member.Id == id {
			return &member
		}
	}
	return nil
}

func LoadMembersFromCsv(membersCsvFile string) []NewForumUser {
	csvInputFile, err := os.Open(membersCsvFile)
	util.Panic(err)
	reader := csv.NewReader(bufio.NewReader(csvInputFile))
	reader.Comma = csvDelimiter

	isHeaderLine := true
	processedLine := 0

	var members []NewForumUser
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
			memberRow := smf_members_row{
				id_member:            line[0],
				member_name:          line[1],
				date_registered:      line[2],
				posts:                line[3],
				id_group:             line[4],
				lngfile:              line[5],
				last_login:           line[6],
				real_name:            line[7],
				instant_messages:     line[8],
				unread_messages:      line[9],
				new_pm:               line[10],
				alerts:               line[11],
				buddy_list:           line[12],
				pm_ignore_list:       line[13],
				pm_prefs:             line[14],
				mod_prefs:            line[15],
				passwd:               line[16],
				email_address:        line[17],
				personal_text:        line[18],
				birthdate:            line[19],
				website_title:        line[20],
				website_url:          line[21],
				show_online:          line[22],
				time_format:          line[23],
				signature:            line[24],
				time_offset:          line[25],
				avatar:               line[26],
				usertitle:            line[27],
				member_ip:            line[28],
				member_ip2:           line[29],
				secret_question:      line[30],
				secret_answer:        line[31],
				id_theme:             line[32],
				is_activated:         line[33],
				validation_code:      line[34],
				id_msg_last_visit:    line[35],
				additional_groups:    line[36],
				smiley_set:           line[37],
				id_post_group:        line[38],
				total_time_logged_in: line[39],
				password_salt:        line[40],
				ignore_boards:        line[41],
				warning:              line[42],
				passwd_flood:         line[43],
				pm_receive_from:      line[44],
				timezone:             line[45],
				tfa_secret:           line[46],
				tfa_backup:           line[47],
				shopMoney:            line[48],
				shopBank:             line[49],
				shopInventory_hide:   line[50],
				gamesPass:            line[51],
				referral:             line[52],
				ref_count:            line[53],
			}
			members = append(members, NewForumUser{
				Id:       memberRow.id_member,
				Username: memberRow.member_name,
				//InventoryRows: nil,
				MemberRow: &memberRow,
				//Items:         nil,
			})
		}
	}
	return members
}

func PairUsersBasedOnCsvId(oldUsers []OldForumUser, newUsers []NewForumUser) []MigratedUser {
	var migratedUsers []MigratedUser
	for _, oldUser := range oldUsers {
		for _, newUser := range newUsers {
			if oldUser.HPHId == newUser.Id {
				fmt.Println(newUser.Id + " - " + oldUser.Id)
				migratedUsers = append(migratedUsers, MigratedUser{
					Username:     oldUser.Username,
					OldForumUser: &oldUser,
					NewForumUser: &newUser,
				})
			}
		}
	}
	return migratedUsers
}
