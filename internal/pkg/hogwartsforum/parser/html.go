package parser

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	htmlpkg "html"
	"localdev/dobby-server/internal/pkg/util"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	gobsSadFace = "https://i.imgur.com/0QvpfXr.gif"
	gobs150     = "https://i.imgur.com/sG6RoCV.gif"
	gobs100     = "https://i.imgur.com/0NIKwhb.gif"
	gobs50      = "https://i.imgur.com/NJ6SzEJ.gif"
	chess200    = "https://i.imgur.com/r7LgPEs.gif"
	chess150    = "https://i.imgur.com/UOSjsVJ.gif"
	chess100    = "https://i.imgur.com/CkHmQYW.gif"
	chess50     = "https://i.imgur.com/XtAGtdR.gif"
)

func GetSubforumAllThreads(html string) []string {
	var threads []string

	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("li.row").Each(func(index int, element *goquery.Selection) {
		text, _ := element.Html()
		threads = append(threads, text)
	})

	return threads
}

func GetSubforumPinnedThreads(html string) []string {
	var pinnedThreads []string

	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Find <li> tags inside <div class="forumbg announcement">
	doc.Find("div.forumbg.announcement li.row").Each(func(index int, element *goquery.Selection) {
		text, _ := element.Html()
		pinnedThreads = append(pinnedThreads, text)
	})

	return pinnedThreads
}

func GetSubforumThreadsNotAnnouncement(html string) []string {
	var threads []string

	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Find <li> tags inside <div class="forumbg"> but not within <div class="forumbg announcement">
	doc.Find("div.forumbg:not(.announcement) li.row").Each(func(index int, element *goquery.Selection) {
		text, _ := element.Html()
		threads = append(threads, text)
	})

	return threads
}

func GetSubforumThreadsAnnouncement(html string) []string {
	var threads []string

	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Find <li> tags inside <div class="forumbg"> but not within <div class="forumbg announcement">
	doc.Find("div.forumbg.announcement li.row").Each(func(index int, element *goquery.Selection) {
		text, _ := element.Html()
		threads = append(threads, text)
	})

	return threads
}

func SubGetThreadUrl(htmlString string) string {
	reader := strings.NewReader(htmlString)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "<Post Title Not Found>"
	}

	// Find the <h2> element with class "topic-title" and get its text
	threadURL := doc.Find("h2.topic-title a").AttrOr("href", "")

	return threadURL
}

func ThreadExtractTitleAndURL(htmlFragment string) (title, url string, err error) {
	reader := strings.NewReader(htmlFragment)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", "", err
	}

	// Find the <a> element inside the <h1> element
	link := doc.Find("h1.page-title a")

	// Extract the title and URL
	title = link.Text()
	url, _ = link.Attr("href")

	return title, url, nil
}

func ThreadExtactReplyData(html string) (tid, t, lt, auth1, auth2 string, err error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	util.Panic(err)

	//find input name tid and get value
	tid, _ = doc.Find("input[name='tid']").Attr("value")

	//find input name t and get value
	t, _ = doc.Find("input[name='t']").Attr("value")

	//find input name lt and get value
	lt, _ = doc.Find("input[name='lt']").Attr("value")

	var secrets []string
	doc.Find("input[name='auth[]']").Each(func(i int, s *goquery.Selection) {
		//get the value of the input
		secret := s.AttrOr("value", "")
		secrets = append(secrets, secret)
	})

	auth1 = ""
	auth2 = ""
	if len(secrets) == 2 {
		auth1 = secrets[0]
		auth2 = secrets[1]
	}

	// if all values are not empty, return them
	if tid != "" && t != "" && lt != "" && auth1 != "" && auth2 != "" {
		return tid, t, lt, auth1, auth2, nil
	} else {
		return "", "", "", "", "", fmt.Errorf("could not extract reply data")
	}

}

func ThreadListPosts(html string) []string {
	var posts []string

	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
		return nil
	}

	doc.Find("div.post").Each(func(index int, element *goquery.Selection) {
		text, _ := element.Html()
		posts = append(posts, text)
	})

	return posts
}

func ThreadNextPageURL(html string) (string, bool) {
	// Load the HTML content into a goquery document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", false
	}

	// Find the <p> tag with the "pagination" class
	paginationElement := doc.Find("p.pagination span")
	if paginationElement.Length() == 0 {
		return "", false // No "next" link found
	}

	nextButton := paginationElement.Find("a").Last()
	nextButtonHtml, _ := nextButton.Html()
	if !strings.Contains(nextButtonHtml, "Siguiente") {
		return "", false // No "next" link found
	}
	nextButtonUrl := nextButton.AttrOr("href", "")
	return nextButtonUrl, true
}

func SubNextPageURL(html string) (string, bool) {
	// Load the HTML content into a goquery document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", false
	}

	// Find the <p> tag with the "pagination" class
	paginationElement := doc.Find("div.pagination span")
	if paginationElement.Length() == 0 {
		return "", false // No "next" link found
	}

	nextButton := paginationElement.Find("a").Last()
	nextButtonHtml, _ := nextButton.Html()
	if !strings.Contains(nextButtonHtml, "Siguiente") {
		return "", false // No "next" link found
	}
	nextButtonUrl := nextButton.AttrOr("href", "")
	return nextButtonUrl, true
}

func PostGetUserName(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "<Username Not Found>"
	}

	post := doc.Find("div.post1")
	username := post.Find("a[href^='/u']").Text()
	return username
}

func PostGetUserUrl(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "<UserSession URL Not Found>"
	}

	post := doc.Find("div.post1")
	userUrl, _ := post.Find("a[href^='/u']").Attr("href")
	return userUrl
}

func PostGetUserHouse(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "<UserSession URL Not Found>"
	}

	information := doc.Find("div.informacion img").Last()
	house := information.AttrOr("alt", "")
	return house
}

func PostGetUrl(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "<Post URL Not Found>"
	}
	var url string
	// Find the <div> element with class "linkfecha"
	linkfecha := doc.Find("div.linkfecha")
	link := linkfecha.Find("a")
	url, _ = link.Attr("href")

	return url
}

func parseDateTime(datetimeStr string) (string, string) {
	// Split the datetime string using '-' or ','
	parts := strings.FieldsFunc(datetimeStr, func(r rune) bool { return r == '-' || r == ',' })

	// Initialize date and time strings
	var dateStr, timeStr string

	if len(parts) == 2 {
		dateStr = parts[0]
		timeStr = strings.TrimSpace(parts[1])
	} else if len(parts) == 3 {
		dateStr = parts[0]
		timeStr = parts[1] + " " + strings.TrimSpace(parts[2])
	}

	return dateStr, timeStr
}

func PostGetDateAndTime(html string, forumDateTime time.Time) *time.Time {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil
	}

	var timeStr, dateStr string

	// Find the <div> element with class "linkfecha" format:  29/11/2023, 08:31
	dateDiv := doc.Find("div.linkfecha").Nodes[0].LastChild
	datetimeStr := strings.TrimSpace(dateDiv.Data)
	if strings.Contains(datetimeStr, "Hoy a las") || strings.Contains(datetimeStr, "Ayer a las") {
		//extract time from "Hoy a las !2:01" or "Ayer a las 12:01"
		timeStr = strings.Split(datetimeStr, " ")[3]
		dateStr = strings.Split(datetimeStr, " ")[0]
		dateStr = util.AdjustDateTimeToStr(forumDateTime, dateStr)
	} else {
		dateStr = strings.Split(datetimeStr, ",")[0]
		timeStr = strings.TrimSpace(strings.Split(datetimeStr, ",")[1])
	}
	layout := "2/1/2006 15:04"
	dateTime, err := time.Parse(layout, dateStr+" "+timeStr)
	util.Panic(err)
	return &dateTime
}

func PostGetEditedDateAndTime(html string) *time.Time {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil
	}

	// Find the <div> element with class "linkfecha"
	dateDiv := doc.Find("div.post3")
	encodedString, _ := dateDiv.Html()

	dateDivStr := htmlpkg.UnescapeString(encodedString)

	if strings.Contains(dateDivStr, "Última edición por") {
		numberIndex := -1
		for i, char := range dateDivStr {
			if char >= '0' && char <= '9' {
				numberIndex = i
				break
			}
		}

		// If a number was found, look for the word ", editado"
		if numberIndex >= 0 {
			editedIndex := strings.Index(dateDivStr, ", editado")
			if editedIndex > numberIndex {
				// Extract the desired part of the string
				result := dateDivStr[numberIndex:editedIndex]
				dateStr := strings.Split(result, ",")[0]
				timeStr := strings.TrimSpace(strings.Split(result, ",")[1])
				layout := "2/1/2006 15:04"
				dateTime, err := time.Parse(layout, dateStr+" "+timeStr)
				util.Panic(err)
				return &dateTime
			}
		}
	}
	return nil
}

func PostGetContent(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "<Post Content Not Found>"
	}
	var content string
	// Find the <div> element with class "content"
	content, _ = doc.Find("div.content").Html()

	return content
}

func PostGetDices(html string) []string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}
	pattern1 := `^Número aleatorio \(\d+,\d+\) : \d+$`
	pattern2 := `^Número aleatorio \(\d+,\d+\) : \(\+\d+\) : \d+$`

	var diceRolls []string
	// Find all <dl> elements with class "codebox"
	doc.Find("dl.codebox").Each(func(i int, dlSelection *goquery.Selection) {
		// Find <dd> tags inside the <dl> element
		ddSelection := dlSelection.Find("dd")
		ddSelection.Each(func(j int, ddSelection *goquery.Selection) {
			// Print the text content of <dd>
			diceLine := ddSelection.Text()
			//match both patterns
			match1, _ := regexp.MatchString(pattern1, diceLine)
			match2, _ := regexp.MatchString(pattern2, diceLine)
			if match1 || match2 {
				diceRolls = append(diceRolls, diceLine)
			}

		})
	})
	return diceRolls
}

func PostGetLinks(contentHtml string) []string {
	reader := strings.NewReader(contentHtml)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	var allLinks []string
	// Find all <a> elements inside the <div class="content">
	doc.Find("a").Each(func(i int, aSelection *goquery.Selection) {
		// Get the href attribute
		link, _ := aSelection.Attr("href")
		allLinks = append(allLinks, link)
	})

	var filteredLinks []string
	for _, link := range allLinks {
		parts := strings.Split(link, "/")
		if len(parts) > 3 {
			firstSegment := parts[3]
			if strings.HasPrefix(firstSegment, "t") { //only threads links
				filteredLinks = append(filteredLinks, link)
			}
		}
	}

	return filteredLinks
}

func IsThreadVisible(html string) bool {
	searchString := `<h1 class="page-title">Informaciones</h1><p>Lo sentimos pero solamente los <strong>usuarios que tengan permisos</strong> pueden leer temas en este foro</p>`

	if strings.Contains(html, searchString) {
		return false
	} else {
		return true
	}
}

func IsLoginCorrect(html string) (bool, string) {
	searchString := "<p>Has escrito un nombre de usuario incorrecto, inactivo o una contraseña inválida."
	searchString2 := "<p>El número máximo de 10 intentos de conexiones autorizadas ha sido superado."

	if strings.Contains(html, searchString) {
		return false, "Has escrito un nombre de usuario incorrecto, inactivo o una contraseña inválida."
	} else if strings.Contains(html, searchString2) {
		return false, "El número máximo de 10 intentos de conexiones autorizadas ha sido superado."
	} else {
		return true, "Inicio de sesión correcto."
	}
}

func GetUsername(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	usernameLine := doc.Find("a.mainmenu").Last().Text()
	username := strings.Split(usernameLine, "[")[1]
	username = strings.Split(username, "]")[0]
	username = strings.TrimSpace(username)
	return username

}

func GetPotionPlayers(html string) (string, string, string, string) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}
	potionInfo := doc.Find("div.xxEDV").Last()
	var player1name, player2name, player1url, player2url string
	potionInfo.Find("strong").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Jugador 1:" {
			player1name = s.NextFiltered("a").Text()
			player1url, _ = s.NextFiltered("a").Attr("href")
		} else if s.Text() == "Jugador 2:" {
			player2name = s.NextFiltered("a").Text()
			player2url, _ = s.NextFiltered("a").Attr("href")
		}
	})
	//remove @ from player names
	player1name = strings.ReplaceAll(player1name, "@", "")
	player2name = strings.ReplaceAll(player2name, "@", "")

	return player1name, player1url, player2name, player2url
}

func GetPotionPlayerProfileUrl(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	player1 := doc.Find("div.postbody").First().Find("strong").First().Text()
	return player1
}

func GetPostSecrets(html string) (string, string) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	//find the input with name "auth[]"
	var secrets []string
	doc.Find("input[name='auth[]']").Each(func(i int, s *goquery.Selection) {
		//get the value of the input
		secret := s.AttrOr("value", "")
		secrets = append(secrets, secret)
	})

	if len(secrets) == 2 {
		return secrets[0], secrets[1]
	} else {
		return "", ""
	}
}

func GetUserTimezone(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	var selectedOptionValue string
	var selectedOptionText string
	//find the <select name="timezone"> element
	doc.Find("select[name='timezone']").Each(func(i int, s *goquery.Selection) {

		//gets value of selected
		selectedOptionValue = s.Find("option[selected]").AttrOr("value", "")
		fmt.Println("selectedOptionValue:", selectedOptionValue)
		//gets text of selected
		selectedOptionText = s.Find("option[selected]").Text()
		fmt.Println("selectedOptionText:", selectedOptionText)

	})

	result := selectedOptionValue + "|" + selectedOptionText
	fmt.Println("Result timezone:", result)

	return result
}

func IsPostSuccessful(html string) (bool, string) {
	searchString := "<p>Tu mensaje ha sido publicado con éxito"

	if strings.Contains(html, searchString) {
		reader := strings.NewReader(html)
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			util.LongPrintlnPrintln("Error:", err)
		}

		//find <a href="/viewtopic?t=91149&amp;topic_name#471117">
		// with text Haz click aquí para ver tu mensaje
		var urlRedirect string
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			if s.Text() == "Haz click aquí para ver tu mensaje" {
				urlRedirect, _ = s.Attr("href")
			}
		})

		//post url is /viewtopic?t=91150&topic_name#471125, extract the value of t
		postUrl := strings.Split(urlRedirect, "t=")[1]
		postUrl = strings.Split(postUrl, "&")[0]
		postUrl = strings.TrimSpace(postUrl)
		postUrl = "t" + postUrl

		return true, urlRedirect
	} else {
		return false, ""
	}
}

func PostGetDateTime(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	util.Panic(err)

	// Find the <div> element with class "linkfecha" format:  29/11/2023, 08:31
	dateDiv := doc.Find("div.linkfecha").Nodes[0].LastChild
	dateDivStr := strings.Trim(dateDiv.Data, " ")
	return dateDivStr
}

func PostIsGobstons(html string) bool {
	//element  <div>
	//    <strong> <User Name> ha efectuado 1 lanzada(s) de uno Gobstons : </strong>
	//    <dl class="codebox"><dd><img src="https://i.imgur.com/0QvpfXr.gif" alt="<User Name> | 01/05/2024 0QvpfXr"/>
	//        </dd>
	//    </dl>
	//</div>
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	searchString := "ha efectuado 1 lanzada(s) de uno Gobstons : "
	// Find the <strong> element with class "lanzada"
	strongElement := doc.Find("dl.codebox")
	haveStrongElement := strongElement.Length() > 0
	haveSearchString := strings.Contains(html, searchString)
	return haveStrongElement && haveSearchString
}

func PostGetGobsValue(html string) int {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	// Find the <img> element inside the <dl> element
	imgElement := doc.Find("dl.codebox img")
	imgSrc, _ := imgElement.Attr("src")
	if imgSrc == gobsSadFace {
		return 0
	} else if imgSrc == gobs50 {
		return 50
	} else if imgSrc == gobs100 {
		return 100
	} else if imgSrc == gobs150 {
		return 150
	} else {
		return -1
	}
}

func PostGetChessValue(html string) int {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	// Find the <img> element inside the <dl> element
	imgElement := doc.Find("dl.codebox img")
	imgSrc, _ := imgElement.Attr("src")
	if imgSrc == chess50 {
		return 50
	} else if imgSrc == chess100 {
		return 100
	} else if imgSrc == chess150 {
		return 150
	} else if imgSrc == chess200 {
		return 200
	} else {
		return -1
	}
}

func PostGetChessLinks(html string) []string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	var allLinks []string
	// Find all <a> elements inside the <div class="content">
	doc.Find("a").Each(func(i int, aSelection *goquery.Selection) {
		// Get the href attribute
		link, _ := aSelection.Attr("href")
		allLinks = append(allLinks, link)
	})

	var filteredLinks []string
	for _, link := range allLinks {
		parts := strings.Split(link, "/")
		if len(parts) > 3 {
			firstSegment := parts[3]
			if strings.HasPrefix(firstSegment, "t") { //only threads links
				filteredLinks = append(filteredLinks, link)
			}
		}
	}

	return filteredLinks
}

type Profile struct {
	Username                   string
	Galeones                   string
	Mensajes                   string
	Ataque                     string
	Defensa                    string
	Edad                       string
	Bando                      string
	Patronus                   string
	Sangre                     string
	Casa                       string
	Inventario1                string
	Inventario2                string
	RazaInventory              Inventory
	HechizosInventory          Inventory
	HabilidadesInventory       Inventory
	HabilidadesDeRazaInventory Inventory
	PocionesInventory          Inventory
	IngredientesInventory      Inventory
	OtrosInventory             Inventory
	LogrosInventory            Inventory
}

func ProfileGetProfile(html string) Profile {
	//create a map <string, string> to store the user profile
	var profile Profile

	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		util.LongPrintlnPrintln("Error:", err)
	}

	doc.Find("h1.page-title").Each(func(i int, s *goquery.Selection) {
		// Get the text inside the <strong> tag
		username := s.Find("strong").Text()
		profile.Username = strings.TrimSpace(username)
	})

	doc.Find("dl#field_id-13").Each(func(i int, s *goquery.Selection) {
		profileField := s.Find("div.field_uneditable").Text()
		profile.Galeones = strings.TrimSpace(profileField)
	})

	doc.Find("dl#field_id-6").Each(func(i int, s *goquery.Selection) {
		profileField := s.Find("div.field_uneditable").Text()
		profile.Mensajes = strings.TrimSpace(profileField)
	})

	doc.Find("dl#field_id9").Each(func(i int, s *goquery.Selection) {
		profileField := s.Find("div.field_uneditable").Text()
		profile.Ataque = strings.TrimSpace(profileField)
	})

	doc.Find("dl#field_id10").Each(func(i int, s *goquery.Selection) {
		profileField := s.Find("div.field_uneditable").Text()
		profile.Defensa = strings.TrimSpace(profileField)
	})

	doc.Find("dl#field_id3").Each(func(i int, s *goquery.Selection) {
		profileField := s.Find("div.field_uneditable").Text()
		profile.Edad = strings.TrimSpace(profileField)
	})

	doc.Find("dl#field_id7").Each(func(i int, s *goquery.Selection) {
		profileField := s.Find("div.field_uneditable").Text()
		profile.Bando = strings.TrimSpace(profileField)
	})

	doc.Find("dl#field_id-11").Each(func(i int, s *goquery.Selection) {
		profileField := s.Find("div.field_uneditable").Text()
		profile.Patronus = strings.TrimSpace(profileField)
	})

	doc.Find("dl#field_id1").Each(func(i int, s *goquery.Selection) {
		profileField := s.Find("div.field_uneditable").Text()
		profile.Sangre = strings.TrimSpace(profileField)
	})

	doc.Find("dl#field_id6").Each(func(i int, s *goquery.Selection) {
		title, exists := s.Find("div.field_uneditable img").Attr("title")
		if exists {
			profile.Casa = strings.TrimSpace(title)
		} else {
			profile.Casa = "ERROR"
		}
	})

	doc.Find("dl#field_id-20").Each(func(i int, s *goquery.Selection) {
		profileField, error := s.Find("div.field_uneditable").Html()
		util.Panic(error)
		profile.Inventario1 = htmlpkg.UnescapeString(profileField)
	})

	doc.Find("dl#field_id11").Each(func(i int, s *goquery.Selection) {
		profileField, error := s.Find("div.field_uneditable").Html()
		util.Panic(error)
		profile.Inventario2 = htmlpkg.UnescapeString(profileField)
	})

	profile.RazaInventory, profile.HechizosInventory, profile.HabilidadesInventory, profile.HabilidadesDeRazaInventory = ExtractItemsBySectionOne(profile.Inventario1)
	profile.PocionesInventory, profile.IngredientesInventory, profile.OtrosInventory, profile.LogrosInventory = ExtractItemsBySectionTwo(profile.Inventario2)

	return profile

}

type Inventory struct {
	Name  string
	Items []ParsedItem
}

type ParsedItem struct {
	Name     string
	Type     string
	ImageUrl string
}

func ExtractItemsBySectionOne(htmlStr string) (Inventory, Inventory, Inventory, Inventory) {
	// Initialize empty sections for each category
	sections := [][]ParsedItem{
		{}, // RAZA (sections[0])
		{}, // HECHIZOS (sections[1])
		{}, // HABILIDADES (sections[2])
		{}, // HABILIDADES DE RAZA (sections[3])
	}

	var currentSectionType string
	var currentSection []ParsedItem

	// Parse the HTML string
	reader := strings.NewReader(htmlStr)
	doc, _ := goquery.NewDocumentFromReader(reader)

	// Iterate over the HTML structure
	doc.Find("strong > i > *").Each(func(i int, s *goquery.Selection) {
		// Check if this is a <div> that starts a new section
		if s.Is("div") {
			// If there's a current section being processed, add it to the sections array
			if len(currentSection) > 0 {
				// Assign currentSection to the appropriate section index based on the currentSectionType
				if currentSectionType == "RAZA" {
					sections[0] = currentSection
				} else if currentSectionType == "HECHIZOS" {
					sections[1] = currentSection
				} else if currentSectionType == "HABILIDADES" {
					sections[2] = currentSection
				} else if strings.HasPrefix(currentSectionType, "HABILIDADES DE") {
					sections[3] = currentSection
				}

				currentSection = []ParsedItem{} // Start a new section
			}
			// Set the type of the current section (e.g., "RAZA", "HECHIZOS")
			currentSectionType = strings.TrimSpace(s.Find("center").Text())
		} else if s.Is("img") {
			// If it's an <img>, create an ParsedItem and add it to the current section
			src, _ := s.Attr("src")
			title, _ := s.Attr("title")
			//replace   with space in the title
			title = strings.ReplaceAll(title, " ", " ")
			item := ParsedItem{
				Name:     strings.TrimSpace(title),
				Type:     currentSectionType,
				ImageUrl: src,
			}
			currentSection = append(currentSection, item)
		}
	})

	// Assign the last section if it has elements
	if len(currentSection) > 0 {
		switch currentSectionType {
		case "RAZA":
			sections[0] = currentSection
		case "HECHIZOS":
			sections[1] = currentSection
		case "HABILIDADES":
			sections[2] = currentSection
		case "HABILIDADES DE RAZA":
			sections[3] = currentSection
		}
	}

	// Create Inventory structs for each section
	razaInventory := Inventory{
		Name:  "RAZA",
		Items: sections[0],
	}
	hechizosInventory := Inventory{
		Name:  "HECHIZOS",
		Items: sections[1],
	}
	habilidadesInventory := Inventory{
		Name:  "HABILIDADES",
		Items: sections[2],
	}
	habilidadesDeRazaInventory := Inventory{
		Name:  "HABILIDADES DE RAZA",
		Items: sections[3],
	}

	return razaInventory, hechizosInventory, habilidadesInventory, habilidadesDeRazaInventory
}

func ExtractItemsBySectionTwo(htmlStr string) (Inventory, Inventory, Inventory, Inventory) {
	// Initialize empty sections for each category
	sections := [][]ParsedItem{
		{}, // POCIONES (sections[0])
		{}, // INGREDIENTES RITUALES (sections[1])
		{}, // OTROS (sections[2])
		{}, // LOGROS (sections[3])
	}

	var currentSectionType string
	var currentSection []ParsedItem

	// Parse the HTML string
	reader := strings.NewReader(htmlStr)
	doc, _ := goquery.NewDocumentFromReader(reader)

	// Iterate over the HTML structure
	doc.Find("strong > i > *").Each(func(i int, s *goquery.Selection) {
		// Check if this is a <div> that starts a new section
		if s.Is("div") {
			// If there's a current section being processed, add it to the sections array
			if len(currentSection) > 0 {
				// Assign currentSection to the appropriate section index based on the currentSectionType
				switch currentSectionType {
				case "POCIONES":
					sections[0] = currentSection
				case "INGREDIENTES RITUALES":
					sections[1] = currentSection
				case "OTROS":
					sections[2] = currentSection
				case "LOGROS":
					sections[3] = currentSection
				}
				currentSection = []ParsedItem{} // Start a new section
			}
			// Set the type of the current section (e.g., "RAZA", "HECHIZOS")
			currentSectionType = strings.TrimSpace(s.Find("center").Text())
		} else if s.Is("img") {
			// If it's an <img>, create an ParsedItem and add it to the current section
			src, _ := s.Attr("src")
			title, _ := s.Attr("title")

			//replace   with space in the title
			title = strings.ReplaceAll(title, " ", " ")
			item := ParsedItem{
				Name:     strings.TrimSpace(title),
				Type:     currentSectionType,
				ImageUrl: src,
			}
			currentSection = append(currentSection, item)
		}
	})

	// Assign the last section if it has elements
	if len(currentSection) > 0 {
		switch currentSectionType {
		case "POCIONES":
			sections[0] = currentSection
		case "INGREDIENTES RITUALES":
			sections[1] = currentSection
		case "OTROS":
			sections[2] = currentSection
		case "LOGROS":
			sections[3] = currentSection
		}
	}

	// Create Inventory structs for each section
	pocionesInventory := Inventory{
		Name:  "POCIONES",
		Items: sections[0],
	}
	ingredientesInventory := Inventory{
		Name:  "INGREDIENTES RITUALES",
		Items: sections[1],
	}
	otrosInventory := Inventory{
		Name:  "OTROS",
		Items: sections[2],
	}
	logrosDeRazaInventory := Inventory{
		Name:  "LOGROS",
		Items: sections[3],
	}

	return pocionesInventory, ingredientesInventory, otrosInventory, logrosDeRazaInventory
}

type ShopCategory struct {
	Name  string
	Items []ShopItem
}
type ShopItem struct {
	Name        string
	ImgUrl      string
	NewImgUrl   string
	Price       string
	Description string
	Shop        string
	Category    string
	Filename    string
}

func ParseShop(htmlStr string, shopName string, listCategories []ShopCategory) []ShopCategory {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	util.Panic(err)

	//var categories []ShopCategory

	// Regular expression to remove special characters
	re := regexp.MustCompile(`[^\w+]`)

	// Iterate over each category
	doc.Find("ul.tabs li span").Each(func(i int, s *goquery.Selection) {
		//categoryName := s.Text()
		//category := ShopCategory{Name: shopName + " - " + categoryName, Items: []ShopItem{}}

		// Find corresponding items in the tab content
		doc.Find("ul.tab__content li").Eq(i).Find("figure").Each(func(j int, f *goquery.Selection) {
			itemName := f.Find("span.i_n").Text()
			itemImgUrl, _ := f.Find("img").Attr("src")
			itemPrice := f.Find("span.nbprix").Text()

			// Replace spaces with '+' and remove special characters
			newImgUrl := strings.Trim(itemName, " ")
			newImgUrl = strings.ReplaceAll(newImgUrl, "á", "a")
			newImgUrl = strings.ReplaceAll(newImgUrl, "é", "e")
			newImgUrl = strings.ReplaceAll(newImgUrl, "í", "i")
			newImgUrl = strings.ReplaceAll(newImgUrl, "ó", "o")
			newImgUrl = strings.ReplaceAll(newImgUrl, "ú", "u")
			newImgUrl = strings.ReplaceAll(newImgUrl, " ", "+")
			newImgUrl = re.ReplaceAllString(newImgUrl, "")

			for i := range listCategories {
				cat := &listCategories[i]

				for j := range cat.Items {
					item := &cat.Items[j]

					if item.Name == strings.Trim(itemName, " ") {
						item.NewImgUrl = newImgUrl
						item.Price = itemPrice
						//item.Shop = shopName
						//item.Category = cat.Name
						item.ImgUrl = itemImgUrl
					}
				}
			}

			/*
				item := ShopItem{
					Name:      itemName,
					ImgUrl:    itemImgUrl,
					NewImgUrl: newImgUrl,
					Price:     itemPrice,
					Shop:      shopName,
					Category:  categoryName,
				}
				category.Items = append(category.Items, item)
			*/
		})

		//categories = append(categories, category)
	})

	return listCategories
}

// FindDescription busca la descripción de un ítem específico en el HTML proporcionado.
func findDescription(descripcionesHtml string, itemName string) string {
	// Cargar el HTML en goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(descripcionesHtml))
	if err != nil {
		fmt.Println("Error al parsear el HTML:", err)
		return ""
	}

	// Buscar el ítem por nombre
	var descriptionBlock string
	doc.Find("div.spoiler_content").Each(func(i int, s *goquery.Selection) {
		s.Find("strong").Each(func(j int, s *goquery.Selection) {
			if s.Text() == itemName {
				descriptionBlock = s.Text() + " " + s.Parent().Text()
				// Eliminar el nombre del ítem del texto completo para obtener solo la descripción
				descriptionBlock = strings.Replace(descriptionBlock, itemName, "", 1)
				descriptionBlock = strings.TrimSpace(descriptionBlock)
			}
		})
	})

	//remove al \n from the descriptionBlock
	descriptionBlock = strings.ReplaceAll(descriptionBlock, "\n", "")
	descriptionBlock = strings.ReplaceAll(descriptionBlock, "\t", "")
	descriptionBlock = formatText(descriptionBlock)

	//split descriptionBlock by "•" and get the first part
	descriptionParts := strings.Split(descriptionBlock, "•")

	if len(descriptionParts) < 1 {
		fmt.Println("Error al obtener la descripción del ítem:", itemName)
	}

	description := ""
	for _, part := range descriptionParts {
		if strings.Contains(part, itemName) {
			description = part
			description = strings.Replace(description, itemName, "", 1)
			description = strings.TrimSpace(description)
			break
		}
	}

	return description
}

// formatText formatea el texto para eliminar múltiples espacios y caracteres no deseados.
func formatText(input string) string {
	re := regexp.MustCompile(`\s+`)
	//formattedText := re.ReplaceAllString(strings.TrimSpace(strings.ReplaceAll(input, "•", "")), " ")
	formattedText := re.ReplaceAllString(strings.TrimSpace(input), " ")

	return formattedText
}

func CreateCategoriesFromDescriptions(htmlStr string) []ShopCategory {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	util.Panic(err)

	var categories []ShopCategory
	var images []string
	var names []string

	// Iterar sobre cada categoría
	doc.Find("dl.codebox.spoiler").Each(func(i int, s *goquery.Selection) {
		var category ShopCategory

		// Obtener el nombre de la categoría
		category.Name = s.Find("dt.spoiler_title").Text()

		// Iterar sobre cada ítem dentro de la categoría
		s.Find("div.spoiler_content").Each(func(j int, item *goquery.Selection) {
			itemHtml, err := item.Html()
			util.Panic(err)

			parts := strings.Split(itemHtml, "<img")
			for _, part := range parts {
				//fmt.Println(part)
				part = formatText(part)
				if part == "" {
					continue
				}
				if part == "<br/>" {
					continue
				}

				if !strings.Contains(part, "strong") || !strings.Contains(part, "br") || !strings.Contains(part, "img") {
					continue
				}

				var shopItem ShopItem
				itemName, err := ExtractStrongText(part)
				util.Panic(err)
				imgUrl, err := ExtractImageURL(part)
				util.Panic(err)
				description, err := ExtractBrText(part)
				util.Panic(err)

				shopItem.Name = strings.TrimSpace(itemName)
				shopItem.ImgUrl = strings.TrimSpace(imgUrl)
				shopItem.Description = strings.TrimSpace(description)
				shopItem.Category = category.Name

				switch shopItem.Category {
				case "Generales":
					shopItem.Shop = "Objetos"
				case "Armas Blancas":
					shopItem.Shop = "Objetos"
				case "Criaturas":
					shopItem.Shop = "Objetos"
				case "Estudiantes":
					shopItem.Shop = "Objetos"
				case "Hogwarts":
					shopItem.Shop = "Objetos"
				case "Ingredientes de Rituales":
					shopItem.Shop = "Objetos"
				case "Ministerio":
					shopItem.Shop = "Objetos"
				case "Mortífagos":
					shopItem.Shop = "Objetos"
				case "Pociones":
					shopItem.Shop = "Objetos"
				case "Propiedades":
					shopItem.Shop = "Objetos"
				case "Quidditch":
					shopItem.Shop = "Objetos"
				case "Transporte Mágico":
					shopItem.Shop = "Objetos"
				case "San Mungo":
					shopItem.Shop = "Objetos"
				case "Autorizaciones":
					shopItem.Shop = "Objetos"
				case "Premios Misiones":
					shopItem.Shop = "Objetos"
				case "Premios Expediciones":
					shopItem.Shop = "Objetos"
				case "Premios Nimbus":
					shopItem.Shop = "Objetos"
				case "Premios Situaciones":
					shopItem.Shop = "Objetos"
				case "Mini-tramas":
					shopItem.Shop = "Objetos"
				case "San Duende":
					shopItem.Shop = "Objetos"
				case "Hechizos Básicos":
					shopItem.Shop = "Hechizos"
				case "Hechizos de Sanación":
					shopItem.Shop = "Hechizos"
				case "Hechizos de Ataque":
					shopItem.Shop = "Hechizos"
				case "Hechizos de Defensa":
					shopItem.Shop = "Hechizos"
				case "Hechizos de Ataque y Defensa":
					shopItem.Shop = "Hechizos"
				case "Hechizos Aurores":
					shopItem.Shop = "Hechizos"
				case "Hechizos Mortífagos":
					shopItem.Shop = "Hechizos"
				case "Maleficios":
					shopItem.Shop = "Hechizos"
				case "Habilidades Adquiribles":
					shopItem.Shop = "Habilidades"
				case "Habilidades de Licántropos":
					shopItem.Shop = "Habilidades"
				case "Habilidades de Semigigantes":
					shopItem.Shop = "Habilidades"
				case "Habilidades Sirenas":
					shopItem.Shop = "Habilidades"
				case "Habilidades de Vampiros":
					shopItem.Shop = "Habilidades"
				case "Habilidades de Veela":
					shopItem.Shop = "Habilidades"
				case "Habilidades de Híbridos Innatas":
					shopItem.Shop = "Habilidades"
				case "Habilidades de Humanos":
					shopItem.Shop = "Habilidades"
				case "Habilidades Adquiridas":
					shopItem.Shop = "Habilidades"
				case "Habilidades Innatas":
					shopItem.Shop = "Habilidades"
				case "Razas":
					shopItem.Shop = "Razas"
				case "Arcana High Bar":
					shopItem.Shop = "Costello"
				case "Borgin & Burkes":
					shopItem.Shop = "Costello"
				case "El Nox":
					shopItem.Shop = "Costello"
				case "Mortem Gemma":
					shopItem.Shop = "Costello"
				case "Portafolio":
					shopItem.Shop = "Costello"
				case "Aquí te tengo tu cariñito":
					shopItem.Shop = "Negocios"
				case "Báthory Square Garden":
					shopItem.Shop = "Negocios"
				case "Botica Slug & Jiggers":
					shopItem.Shop = "Negocios"
				case "Chez Winnie":
					shopItem.Shop = "Negocios"
				case "Danceteria Rolling Hall":
					shopItem.Shop = "Negocios"
				case "El Lux":
					shopItem.Shop = "Negocios"
				case "FLEUR":
					shopItem.Shop = "Negocios"
				case "Flourish & Blotts":
					shopItem.Shop = "Negocios"
				case "Luxxuria":
					shopItem.Shop = "Negocios"
				case "Moonlight Shadow Planetary":
					shopItem.Shop = "Negocios"
				case "Mystic Momentum":
					shopItem.Shop = "Negocios"
				case "Nym's Treasure":
					shopItem.Shop = "Negocios"
				case "Peonie's Ribbon":
					shopItem.Shop = "Negocios"
				case "Plants & Seeds":
					shopItem.Shop = "Negocios"
				case "Ragnarok":
					shopItem.Shop = "Negocios"
				case "Rose":
					shopItem.Shop = "Negocios"
				case "Royal Vauxhall Tavern":
					shopItem.Shop = "Negocios"
				case "Sacred Lotus":
					shopItem.Shop = "Negocios"
				case "Saint Ellis Hospital":
					shopItem.Shop = "Negocios"
				case "Sortilegios Weasley":
					shopItem.Shop = "Negocios"
				case "Tierra y Cristal":
					shopItem.Shop = "Negocios"
				case "Vinos Zabini":
					shopItem.Shop = "Negocios"
				case "Wizarding World Curse":
					shopItem.Shop = "Negocios"
				case "Logros de Colaborador del Mes":
					shopItem.Shop = "Logros"
				case "Logros de Personaje del Mes":
					shopItem.Shop = "Logros"
				case "Logros de Awards":
					shopItem.Shop = "Logros"
				case "Logros de Duelos":
					shopItem.Shop = "Logros"
				case "Logros de Misiones":
					shopItem.Shop = "Logros"
				case "Logros de Expediciones":
					shopItem.Shop = "Logros"
				case "Logros de Pociones":
					shopItem.Shop = "Logros"
				case "Logros de Nimbus":
					shopItem.Shop = "Logros"
				case "Logros de Situaciones":
					shopItem.Shop = "Logros"
				case "Logros de Cámara de Creación Mágica":
					shopItem.Shop = "Logros"
				case "Logros de Rituales":
					shopItem.Shop = "Logros"
				case "Logros de San Mungo":
					shopItem.Shop = "Logros"
				case "Logros de Costello":
					shopItem.Shop = "Logros"
				case "Logros de Colección de Cromos":
					shopItem.Shop = "Logros"
				case "Logros de Top Posteadores":
					shopItem.Shop = "Logros"
				case "Logros de Empleo":
					shopItem.Shop = "Logros"
				case "Logros de Hall of Fame":
					shopItem.Shop = "Logros"
				case "Logros de Torneo de Duelos":
					shopItem.Shop = "Logros"
				case "Logros de Torneo de Pociones":
					shopItem.Shop = "Logros"
				case "Logros de Torneo de Quidditch Libre":
					shopItem.Shop = "Logros"
				case "Logros de Desafío Hogwarts":
					shopItem.Shop = "Logros"
				case "Logros de Mortífagos":
					shopItem.Shop = "Logros"
				case "Logros de Aurores":
					shopItem.Shop = "Logros"
				}

				//insert image into images array if it is not already there
				if !util.Contains(images, shopItem.ImgUrl) {
					images = append(images, shopItem.ImgUrl)
				} else {
					fmt.Println(shopItem.Name+": Image already in array: ", shopItem.ImgUrl)
				}
				// insert name into names array if it is not already there
				if !util.Contains(names, shopItem.Name) {
					names = append(names, shopItem.Name)
				} else {
					fmt.Println(shopItem.Name + ": Name already in array")
				}

				if shopItem.Name != "aquí" {
					category.Items = append(category.Items, shopItem)
				}
			}
		})

		categories = append(categories, category)
	})

	return categories
}

func ExtractStrongText(input string) (string, error) {
	startTag := "<strong>"
	endTag := "</strong>"

	// Buscar el índice del inicio de la etiqueta <strong>
	startIndex := strings.Index(input, startTag)
	if startIndex == -1 {
		return "", errors.New("no se encontró la etiqueta de inicio <strong> en " + input)
	}

	// Ajustar el índice para empezar justo después de <strong>
	startIndex += len(startTag)

	// Buscar el índice del final de la etiqueta </strong> después del <strong>
	endIndex := strings.Index(input[startIndex:], endTag)
	if endIndex == -1 {
		return "", errors.New("no se encontró la etiqueta de cierre </strong>" + input)
	}

	// Extraer el texto entre las etiquetas <strong> y </strong>
	text := input[startIndex : startIndex+endIndex]
	return text, nil
}

func ExtractBrText(input string) (string, error) {
	startTag := "<br/>"
	endTag := "<br/>"

	// Buscar el índice del inicio de la etiqueta <strong>
	startIndex := strings.Index(input, startTag)
	if startIndex == -1 {
		return "", errors.New("no se encontró la etiqueta de inicio <br/> en " + input)
	}

	// Ajustar el índice para empezar justo después de <strong>
	startIndex += len(startTag)

	// Buscar el índice del final de la etiqueta </strong> después del <strong>
	endIndex := strings.Index(input[startIndex:], endTag)
	if endIndex == -1 {
		return "", errors.New("no se encontró la etiqueta de cierre <br/> en " + input)
	}

	// Extraer el texto entre las etiquetas <strong> y </strong>
	text := input[startIndex : startIndex+endIndex]
	return text, nil
}

func ExtractImageURL(input string) (string, error) {
	// La cadena que buscamos dentro del input es `src="`
	srcPrefix := `src="`
	startIndex := strings.Index(input, srcPrefix)
	if startIndex == -1 {
		return "", errors.New("no se encontró la cadena 'src=\"'")
	}

	// Ajustar el índice para que apunte al inicio de la URL
	startIndex += len(srcPrefix)

	// Buscar el índice del final de la URL, que está delimitado por `"`
	endIndex := strings.Index(input[startIndex:], `"`)
	if endIndex == -1 {
		return "", errors.New("no se encontró el final de la URL (delimitado por comillas)")
	}

	// Extraer la URL de la imagen
	url := input[startIndex : startIndex+endIndex]
	return url, nil
}
