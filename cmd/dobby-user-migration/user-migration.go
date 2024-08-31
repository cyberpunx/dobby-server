package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	loginUsername = "Desarrollo"
	loginPassword = "programación2055"
	//UsersToMigrateFromOldForum = "Datos de usuarios - Migrados HR.csv"
	UsersToMigrateFromOldForum = "TEST.csv"
	ItemTableOnNewForum        = "smf_stshop_items.csv"
	ItemUrlsCsv                = "items.csv"
	MemberTableOnNewForum      = "smf_members.csv"
	InventoryTable             = "smf_stshop_inventory.csv" +
		""
	csvDelimiter = ','
)

type session struct {
	Tool *tool.Tool
	Conf config
}

type config struct {
	BaseUrl string `json:"baseUrl"`
}

var CurrentUser = ""

type Item struct {
	Itemid           string
	name             string
	image            string
	description      string
	price            string
	stock            string
	module           string
	info1            string
	info2            string
	info3            string
	info4            string
	input_needed     string
	can_use_item     string
	delete_after_use string
	catid            string
	status           string
	itemlimit        string
	imgurUrl         string
}

type OldForumUser struct {
	Id         string
	HPHId      string
	Username   string
	Profile    parser.Profile
	TotalItems []parser.ParsedItem
}

type NewForumUser struct {
	Id            string
	Username      string
	InventoryRows *[]smf_shop_inventory_row
	MemberRow     *smf_members_row
	CustomFields  *[]smf_themes
	Items         *[]Item
	NotFoundItems []parser.ParsedItem
}

type MigratedUser struct {
	Username     string
	OldForumUser *OldForumUser
	NewForumUser *NewForumUser
}

type smf_themes struct {
	id_member string
	id_theme  string
	variable  string
	value     string
}

type smf_shop_inventory_row struct {
	userid    string
	itemid    string
	trading   string
	tradecost string
	date      string
	tradedate string
	fav       string
	name      string
}

// "id_member","member_name","date_registered","posts","id_group","lngfile","last_login","real_name","instant_messages","unread_messages","new_pm","alerts","buddy_list","pm_ignore_list","pm_prefs","mod_prefs","passwd","email_address","personal_text","birthdate","website_title","website_url","show_online","time_format","signature","time_offset","avatar","usertitle","member_ip","member_ip2","secret_question","secret_answer","id_theme","is_activated","validation_code","id_msg_last_visit","additional_groups","smiley_set","id_post_group","total_time_logged_in","password_salt","ignore_boards","warning","passwd_flood","pm_receive_from","timezone","tfa_secret","tfa_backup","shopMoney","shopBank","shopInventory_hide","gamesPass","referral","ref_count"
type smf_members_row struct {
	id_member            string
	member_name          string
	date_registered      string
	posts                string
	id_group             string
	lngfile              string
	last_login           string
	real_name            string
	instant_messages     string
	unread_messages      string
	new_pm               string
	alerts               string
	buddy_list           string
	pm_ignore_list       string
	pm_prefs             string
	mod_prefs            string
	passwd               string
	email_address        string
	personal_text        string
	birthdate            string
	website_title        string
	website_url          string
	show_online          string
	time_format          string
	signature            string
	time_offset          string
	avatar               string
	usertitle            string
	member_ip            string
	member_ip2           string
	secret_question      string
	secret_answer        string
	id_theme             string
	is_activated         string
	validation_code      string
	id_msg_last_visit    string
	additional_groups    string
	smiley_set           string
	id_post_group        string
	total_time_logged_in string
	password_salt        string
	ignore_boards        string
	warning              string
	passwd_flood         string
	pm_receive_from      string
	timezone             string
	tfa_secret           string
	tfa_backup           string
	shopMoney            string
	shopBank             string
	shopInventory_hide   string
	gamesPass            string
	referral             string
	ref_count            string
}

func main() {
	//LOADS ITEMS FROM CSV
	itemList := LoadItemsFromCsv()
	fmt.Println("Cantidad de items cargados: ", len(*itemList))

	// FORUM LOGIN
	s := forumLogin(loginUsername, loginPassword)
	oldUsers := LoadOldForumUsersFromCsv(UsersToMigrateFromOldForum, s)
	newUsers := LoadMembersFromCsv(MemberTableOnNewForum)

	//PAIR USERS BASED ON USERNAME
	migratedUsers := PairUsersBasedOnCsvId(oldUsers, newUsers)
	PairItemsToUsers(migratedUsers, itemList)

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

func LoadItemsFromCsv() *[]Item {
	csvInputFile, err := os.Open(ItemTableOnNewForum)
	util.Panic(err)

	reader := csv.NewReader(bufio.NewReader(csvInputFile))
	reader.Comma = csvDelimiter

	isHeaderLine := true
	processedLine := 0

	var items []Item
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
			itemid := line[0]
			name := line[1]
			image := line[2]
			description := line[3]
			price := line[4]
			stock := line[5]
			module := line[6]
			info1 := line[7]
			info2 := line[8]
			info3 := line[9]
			info4 := line[10]
			input_needed := line[11]
			can_use_item := line[12]
			delete_after_use := line[13]
			catid := line[14]
			status := line[15]
			itemlimit := line[16]

			items = append(items, Item{
				Itemid:           itemid,
				name:             name,
				image:            image,
				description:      description,
				price:            price,
				stock:            stock,
				module:           module,
				info1:            info1,
				info2:            info2,
				info3:            info3,
				info4:            info4,
				input_needed:     input_needed,
				can_use_item:     can_use_item,
				delete_after_use: delete_after_use,
				catid:            catid,
				status:           status,
				itemlimit:        itemlimit,
			})
		}
	}

	for i, item := range items {
		imgurIrl := SearchItemUrlOnCsv(item.name)
		if imgurIrl == "NOT FOUND" {
			fmt.Println("No se encontró la url de la imagen para el item: ", item.name)
		}
		items[i].imgurUrl = imgurIrl
	}

	return &items
}

func SearchItemUrlOnCsv(itemName string) string {
	csvUrlsFile, err := os.Open(ItemUrlsCsv)
	util.Panic(err)
	readerUrls := csv.NewReader(bufio.NewReader(csvUrlsFile))
	readerUrls.Comma = '|' // CSV uses semicolon as delimiter

	isHeaderLine := true
	processedLine := 0

	for {
		line, err := readerUrls.Read()
		if err == io.EOF {
			break
		}
		util.Panic(err)
		if isHeaderLine {
			isHeaderLine = false
		} else {
			processedLine++
			if line[0] == itemName {
				return line[1]
			}
		}
	}
	return "NOT FOUND"
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
				Id:            memberRow.id_member,
				Username:      memberRow.member_name,
				InventoryRows: nil,
				MemberRow:     &memberRow,
				Items:         nil,
			})
		}
	}
	return members
}

func LoadInventoryTableFromCsv(inventoryCsvFile string) []smf_shop_inventory_row {
	csvInputFile, err := os.Open(inventoryCsvFile)
	util.Panic(err)
	reader := csv.NewReader(bufio.NewReader(csvInputFile))
	reader.Comma = csvDelimiter

	isHeaderLine := true
	processedLine := 0

	var inventoryRows []smf_shop_inventory_row
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
			inventoryRow := smf_shop_inventory_row{
				userid:    line[1],
				itemid:    line[2],
				trading:   line[3],
				tradecost: line[4],
				date:      line[5],
				tradedate: line[6],
				fav:       line[7],
			}
			inventoryRows = append(inventoryRows, inventoryRow)
		}
	}
	return inventoryRows
}

func PairUsersBasedOnUsername(oldUsers []OldForumUser, newUsers []NewForumUser) []MigratedUser {
	var migratedUsers []MigratedUser
	for _, oldUser := range oldUsers {
		for _, newUser := range newUsers {
			if oldUser.Username == newUser.Username {
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

func SearchItemByName(itemName string, items []Item) *Item {
	//Catch some badly written itemNames
	if itemName == "Solución Encongedora" {
		itemName = "Solución Encogedora"
	}
	if strings.HasPrefix(itemName, "Autorización Premium") {
		itemName = "Autorización Premium"
	}
	if strings.HasPrefix(itemName, "Autorización Luxury") {
		itemName = "Autorización Luxury"
	}
	if strings.HasPrefix(itemName, "Autorización Standard") {
		itemName = "Autorización Standard"
	}
	if strings.HasPrefix(itemName, "PJ del Mes") {
		itemName = "PJ del Mes"
	}
	if strings.HasPrefix(itemName, "Award") {
		itemName = "Award"
	}
	if itemName == "Maté un Personaje" {
		itemName = "Maté a un Personaje"
	}
	if itemName == "Máscara 1" {
		itemName = "Máscara de Mortífago 1"
	}
	if itemName == "Máscara 2" {
		itemName = "Máscara de Mortífago 2"
	}
	if itemName == "Máscara 3" {
		itemName = "Máscara de Mortífago 3"
	}
	if itemName == "Máscara 4" {
		itemName = "Máscara de Mortífago 4"
	}
	if itemName == "Máscara 5" {
		itemName = "Máscara de Mortífago 5"
	}
	if itemName == "Máscara 6" {
		itemName = "Máscara de Mortífago 6"
	}
	if itemName == "Máscara 7" {
		itemName = "Máscara de Mortífago 7"
	}
	if itemName == "Máscara 8" {
		itemName = "Máscara de Mortífago 8"
	}
	if itemName == "Máscara 9" {
		itemName = "Máscara de Mortífago 9"
	}
	if itemName == "Máscara 10" {
		itemName = "Máscara de Mortífago 10"
	}
	if itemName == "Lechuza Maori" {
		itemName = "Lechuza Maorí"
	}
	if itemName == "Mimempha" {
		itemName = "Mimempha +3"
	}
	if itemName == "Protego Totallum" {
		itemName = "Protego Totallum +3"
	}
	if itemName == "Protego Totallum" {
		itemName = "Protego Totallum +3"
	}
	if itemName == "Protego Totallum" {
		itemName = "Protego Totallum +3"
	}
	if itemName == "Chrono Exsecutio" {
		itemName = "Chrono Exsecutio +3"
	}
	if itemName == "Restex" {
		itemName = "Restex +3"
	}
	if itemName == "Nebulae Catenis" {
		itemName = "Nebulae Catenis +3"
	}
	if itemName == "Fortificum" {
		itemName = "Fortificum  +3"
	}
	if itemName == "Moenia Crystalis" {
		itemName = "Moenia Crystalis +3"
	}
	if itemName == "Chrono Exsecutio" {
		itemName = "Chrono Exsecutio +3"
	}
	if itemName == "Inferi" {
		itemName = "Inferi +3"
	}
	if itemName == "Tabú" {
		itemName = "Tabu"
	}
	if itemName == "Seccionatus" {
		itemName = "Seccionatus +3"
	}

	//remove consecutive spaces from itemName
	itemName = strings.Join(strings.Fields(itemName), " ")
	//trim spaces from itemName
	itemName = strings.TrimSpace(itemName)
	for _, item := range items {
		if strings.ToLower(item.name) == strings.ToLower(itemName) {
			return &item
		}
	}
	return nil
}

func SearchItemByImguUrl(imgurUrl string, items []Item) *Item {
	for _, item := range items {
		if item.imgurUrl == imgurUrl {
			return &item
		}
	}
	return nil
}

func PairItemsToUsers(migratedUsers []MigratedUser, items *[]Item) {
	var totalFoundItems []Item
	var totalNotFoundItems []parser.ParsedItem

	for _, migratedUser := range migratedUsers {
		var userFoundItems []Item
		var userNotFoundItems []parser.ParsedItem

		for _, oldItem := range migratedUser.OldForumUser.TotalItems {
			foundItem := SearchItem(oldItem, *items)
			if foundItem != nil {
				userFoundItems = append(userFoundItems, *foundItem)
			} else {
				userNotFoundItems = append(userNotFoundItems, oldItem)
			}
		}

		if len(userNotFoundItems) > 0 {
			CreateItems(userNotFoundItems, migratedUser.Username)
			fmt.Println(migratedUser.Username+" items para insertar: ", len(userNotFoundItems))
		}

		migratedUser.NewForumUser.Items = &userFoundItems
		migratedUser.NewForumUser.NotFoundItems = userNotFoundItems

		smfThemeRowAttack := smf_themes{
			id_member: migratedUser.NewForumUser.Id,
			id_theme:  "1",
			variable:  "cust_ataque",
			value:     migratedUser.OldForumUser.Profile.Ataque,
		}

		smfThemeRowDefense := smf_themes{
			id_member: migratedUser.NewForumUser.Id,
			id_theme:  "1",
			variable:  "cust_defens",
			value:     migratedUser.OldForumUser.Profile.Defensa,
		}

		migratedUser.NewForumUser.CustomFields = &[]smf_themes{smfThemeRowAttack, smfThemeRowDefense}

		totalFoundItems = append(totalFoundItems, userFoundItems...)
		totalNotFoundItems = append(totalNotFoundItems, userNotFoundItems...)

		migratedUser.NewForumUser.InventoryRows = CreateInventoryRows(userFoundItems, migratedUser.NewForumUser.Id)
		WriteInventoryQuery(*migratedUser.NewForumUser.InventoryRows, migratedUser.Username)
		WriteUpdateMemberQuery(&migratedUser)
		WriteUpdateThemesQuery(*migratedUser.NewForumUser.CustomFields, migratedUser.Username)
		//WriteQueryAddRemainingInventoryItems(migratedUser)
	}

	//itemsToInsertIntoDatabase are the totalNotFoundItems without the duplicates
	//an item is a duplicate if the item Name and the item ImageUrl are the same
	var itemsToInsertIntoDatabase []parser.ParsedItem

	for _, item := range totalNotFoundItems {
		isDuplicate := false
		for _, itemToInsert := range itemsToInsertIntoDatabase {
			if item.Name == itemToInsert.Name && item.ImageUrl == itemToInsert.ImageUrl {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			itemsToInsertIntoDatabase = append(itemsToInsertIntoDatabase, item)
		}
	}

	fmt.Println("Total de items encontrados: ", len(totalFoundItems))
	fmt.Println("Total de items no encontrados: ", len(totalNotFoundItems))
	fmt.Println("Total de items a insertar en la base de datos: ", len(itemsToInsertIntoDatabase))

	//Create SQL quety to insert items into database
	CreateItems(itemsToInsertIntoDatabase, "ITEMS-SIN-CATEGORIA")

	//ReadReport()
}

func WriteQueryAddRemainingInventoryItems(migratedUser MigratedUser) {
	inventoryRows := LoadInventoryTableFromCsv(InventoryTable)

	var remainingItemsToInsertRows []smf_shop_inventory_row

	for _, item := range *migratedUser.NewForumUser.Items {
		foundRow := false
		for _, row := range inventoryRows {
			if row.userid == migratedUser.NewForumUser.Id && row.itemid == item.Itemid {
				foundRow = true
				break
			}
		}
		if !foundRow {
			remainingItemsToInsertRows = append(remainingItemsToInsertRows, smf_shop_inventory_row{
				userid:    migratedUser.NewForumUser.Id,
				itemid:    item.Itemid,
				trading:   "0",
				tradecost: "0",
				date:      "0",
				tradedate: "0",
				fav:       "0",
				name:      item.name,
			})
		}
	}

	sqlQuery := ""
	sqlQuery += "-- " + migratedUser.Username + "\n"
	for _, row := range remainingItemsToInsertRows {
		sqlQuery += fmt.Sprintf("INSERT INTO `smf_stshop_inventory` (`userid`, `itemid`, `trading`, `tradecost`, `date`, `tradedate`, `fav`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s'); --%s", row.userid, row.itemid, row.trading, row.tradecost, row.date, row.tradedate, row.fav, row.name)
		sqlQuery += "\n"
	}
	sqlQuery += "\n\n"

	//append to file if exists
	f, err := os.OpenFile("add remaining items", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//append total_Lines to file
	_, err = f.WriteString(sqlQuery)
	util.Panic(err)
}

func WriteUpdateThemesQuery(customFields []smf_themes, username string) {
	//update or insert into smf_themes, columns cust_attack and cust_defense
	sqlQuery := ""
	sqlQuery += "-- " + username + "\n"
	for _, customField := range customFields {
		sqlQuery += fmt.Sprintf("INSERT INTO `smf_themes` (`id_member`, `id_theme`, `variable`, `value`) VALUES ('%s', '%s', '%s', '%s') ON DUPLICATE KEY UPDATE `value` = '%s';", customField.id_member, customField.id_theme, customField.variable, customField.value, customField.value)
	}
	sqlQuery += "\n\n"

	//append to file if exists
	f, err := os.OpenFile("update themes", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//append total_Lines to file
	_, err = f.WriteString(sqlQuery)
	util.Panic(err)
}

func WriteUpdateMemberQuery(user *MigratedUser) {
	//update smf_members, columns shopMoney and posts onl
	oldPosts, err := strconv.Atoi(user.OldForumUser.Profile.Mensajes)
	util.Panic(err)
	newPosts, err := strconv.Atoi(user.NewForumUser.MemberRow.posts)
	util.Panic(err)

	oldMoney, err := strconv.Atoi(user.OldForumUser.Profile.Galeones)
	util.Panic(err)
	newMoney, err := strconv.Atoi(user.NewForumUser.MemberRow.shopMoney)
	util.Panic(err)

	totalPosts := strconv.Itoa(oldPosts + newPosts)
	totalMoney := strconv.Itoa(oldMoney + newMoney)

	sqlQuery := ""
	sqlQuery += "-- " + user.NewForumUser.MemberRow.member_name + "\n"
	sqlQuery += fmt.Sprintf("UPDATE `smf_members` SET `shopMoney` = '%s', `posts` = '%s' WHERE `id_member` = '%s';", totalMoney, totalPosts, user.NewForumUser.Id)
	sqlQuery += "\n\n"

	//append to file if exists
	f, err := os.OpenFile("update money and posts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//append total_Lines to file
	_, err = f.WriteString(sqlQuery)
	util.Panic(err)

}

func WriteInventoryQuery(rows []smf_shop_inventory_row, username string) {
	sqlQuery := ""
	sqlQuery += "-- " + username + "\n"
	for _, row := range rows {
		userId := row.userid
		itemId := row.itemid
		trading := row.trading
		tradecost := row.tradecost
		date := row.date
		tradedate := row.tradedate
		fav := row.fav
		sqlQuery += fmt.Sprintf("INSERT INTO `smf_stshop_inventory` (`userid`, `itemid`, `trading`, `tradecost`, `date`, `tradedate`, `fav`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');", userId, itemId, trading, tradecost, date, tradedate, fav)
		sqlQuery += "\n"
	}
	sqlQuery += "\n\n"

	//append to file if exists
	f, err := os.OpenFile("insert inventory", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//append total_Lines to file
	_, err = f.WriteString(sqlQuery)
	util.Panic(err)
}

func CreateInventoryRows(items []Item, userId string) *[]smf_shop_inventory_row {
	var inventoryRows []smf_shop_inventory_row
	for _, item := range items {
		inventoryRows = append(inventoryRows, smf_shop_inventory_row{
			userid:    userId,
			itemid:    item.Itemid,
			trading:   "0",
			tradecost: "0",
			date:      "0",
			tradedate: "0",
			fav:       "0",
		})
	}
	return &inventoryRows
}

func ReadReport() {
	data, err := ioutil.ReadFile("items_to_insert.html")
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}

	// Convierte el contenido a string
	content := string(data)
	re := regexp.MustCompile(`<strong>(.*?)<\/strong>`)
	matches := re.FindAllStringSubmatch(content, -1)

	// Mapa para contar las ocurrencias
	counts := make(map[string]int)

	// Itera sobre los matches y cuenta las ocurrencias
	for _, match := range matches {
		item := strings.TrimSpace(match[1])
		counts[item]++
	}

	// Convertir el mapa a un slice de pares (item, count)
	type kv struct {
		Key   string
		Value int
	}

	var sortedItems []kv
	for k, v := range counts {
		sortedItems = append(sortedItems, kv{k, v})
	}

	// Ordenar el slice por el valor (count) de mayor a menor
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].Value > sortedItems[j].Value
	})

	// Imprimir los resultados ordenados
	for _, kv := range sortedItems {
		fmt.Printf("%s: %d\n", kv.Key, kv.Value)
	}

}

func CreateItems(itemListToCreate []parser.ParsedItem, username string) {
	total_Lines := "<!-- " + username + " -->\n"
	total_Lines += "<div class=\"spoiler_content\">\n\n"
	for _, itemToCreate := range itemListToCreate {
		imgurUrl := itemToCreate.ImageUrl
		name := itemToCreate.Name

		if name == "" {
			fmt.Println("Nombre vacío: ")
			continue
		}
		if imgurUrl == "" {
			fmt.Println("Url de imagen vacía")
			continue
		}

		descLine := "\t<img src=\"{url}\"/>\n\t• <strong>{Nombre}</strong><br />{descripcion}<br />"
		descLine = strings.ReplaceAll(descLine, "{url}", imgurUrl)
		descLine = strings.ReplaceAll(descLine, "{Nombre}", name)
		descLine = strings.ReplaceAll(descLine, "{descripcion}", "")
		total_Lines += descLine + "\n\n"
	}
	total_Lines += "</div>\n\n"

	//append to file if exists
	f, err := os.OpenFile("items_to_insert.html", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//append total_Lines to file
	_, err = f.WriteString(total_Lines)

	util.Panic(err)
}

func SearchItem(searchItem parser.ParsedItem, items []Item) *Item {
	foundItem := SearchItemByName(searchItem.Name, items)
	if foundItem != nil {
		return foundItem
	} else {
		foundItem = SearchItemByImguUrl(searchItem.ImageUrl, items)
		if foundItem != nil {
			return foundItem
		}
	}

	return nil
}

func replaceHtmlChars(input string) string {
	replacements := map[string]string{
		"&#34;": "\"",
		"&#39;": "'",
		"&lt;":  "<",
		"&gt;":  ">",
		"&amp;": "&",
		// Agrega más reemplazos según sea necesario
	}

	for htmlChar, replacement := range replacements {
		input = strings.ReplaceAll(input, htmlChar, replacement)
	}

	return input
}
