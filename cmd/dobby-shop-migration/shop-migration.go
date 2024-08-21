package main

import (
	"fmt"
	"io"
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	loginUsername = "Desarrollo"
	loginPassword = "programación2055"
	csvDelimiter  = ','

	costelloInput     = "Costello"
	habilidadesInput  = "Habilidades"
	hechizosInput     = "Hechizos"
	negociosInput     = "Negocios"
	objetosInput      = "Objetos"
	descricionesInput = "descripciones.html"

	outputDir      = "output"
	imageOutputDir = outputDir + "/migrated"
	sqlOutputFile  = outputDir + "/insert items.sql"
	csvOutputFile  = outputDir + "/items.csv"
)

var Shop = ""

type session struct {
	Tool *tool.Tool
	Conf config
}

type config struct {
	BaseUrl string `json:"baseUrl"`
}

func main() {
	// FORUM LOGIN
	/*
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

	*/
	descripcionesContents, err := os.ReadFile(descricionesInput)
	util.Panic(err)
	descripcionesHtmlStr := string(descripcionesContents)
	categories := parser.CreateCategoriesFromDescriptions(descripcionesHtmlStr)
	//fmt.Println(fmt.Sprintf("%s\n", util.MarshalJsonPretty(categories)))

	//array with inputs
	inputs := []string{costelloInput, habilidadesInput, hechizosInput, negociosInput, objetosInput}
	for _, input := range inputs {
		Shop = input

		shopContents, err := os.ReadFile(Shop + ".html")
		util.Panic(err)
		shopHtmlStr := string(shopContents)

		categories = parser.ParseShop(shopHtmlStr, Shop, categories)
		//fmt.Println(fmt.Sprintf("%s\n", util.MarshalJsonPretty(shops)))
		/*
			for _, cat := range categories {
				err = generateUrlList(cat)
				util.Panic(err)
			}

		*/

		for _, cat := range categories {
			err = processImages(&cat)
			util.Panic(err)
		}
	}

	//fmt.Println(fmt.Sprintf("%s\n", util.MarshalJsonPretty(categories)))

	/*
		var noShopItems []parser.ShopItem
		for i, _ := range categories {
			cat := &categories[i]
			for j, _ := range cat.Items {
				item := &cat.Items[j]


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



				if item.Shop == "" {
					noShopItems = append(noShopItems, *item)
				}

				if len(noShopItems) > 0 {
					fmt.Println("!!! Items sin tienda:")
					fmt.Println(fmt.Sprintf("%s\n", util.MarshalJsonPretty(noShopItems)))
				}
			}
		}

	*/

	queries := generateSqlQuery(categories)
	//create query file
	err = os.WriteFile(sqlOutputFile, []byte(queries), 0644)
	util.Panic(err)

	//fmt.Println(queries)
}

type smf_stshop_items struct {
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

func generateSqlQuery(categories []parser.ShopCategory) string {
	catList := map[string]int{
		"Objetos - Generales":                           1,
		"Objetos - Armas Blancas":                       2,
		"Objetos - Criaturas":                           3,
		"Objetos - Estudiantes":                         4,
		"Objetos - Hogwarts":                            5,
		"Objetos - Ingredientes de Rituales":            6,
		"Objetos - Ministerio":                          7,
		"Objetos - Mortífagos":                          8,
		"Objetos - Pociones":                            9,
		"Objetos - Propiedades":                         10,
		"Objetos - Quidditch":                           11,
		"Objetos - Transporte Mágico":                   12,
		"Objetos - San Mungo":                           13,
		"Objetos - Autorizaciones":                      14,
		"Objetos - Premios Misiones":                    15,
		"Objetos - Premios Expediciones":                16,
		"Objetos - Premios Nimbus":                      17,
		"Objetos - Premios Situaciones":                 18,
		"Objetos - Mini-tramas":                         19,
		"Objetos - San Duende":                          20,
		"Hechizos - Hechizos Básicos":                   21,
		"Hechizos - Hechizos de Sanación":               22,
		"Hechizos - Hechizos de Ataque":                 23,
		"Hechizos - Hechizos de Defensa":                24,
		"Hechizos - Hechizos de Ataque y Defensa":       25,
		"Hechizos - Hechizos Aurores":                   26,
		"Hechizos - Hechizos Mortífagos":                27,
		"Hechizos - Maleficios":                         28,
		"Razas - Razas":                                 29,
		"Habilidades - Habilidades Innatas":             30,
		"Habilidades - Habilidades Adquiribles":         31,
		"Habilidades - Habilidades de Licántropos":      32,
		"Habilidades - Habilidades de Semigigantes":     33,
		"Habilidades - Habilidades Sirenas":             34,
		"Habilidades - Habilidades de Vampiros":         35,
		"Habilidades - Habilidades de Veela":            36,
		"Habilidades - Habilidades de Híbridos Innatas": 37,
		"Habilidades - Habilidades de Humanos":          38,
		"Costello - Arcana High Bar":                    39,
		"Costello - Borgin & Burkes":                    40,
		"Costello - El Nox":                             41,
		"Costello - Mortem Gemma":                       42,
		"Costello - Portafolio":                         43,
		"Negocios - Aquí te tengo tu cariñito":          44,
		"Negocios - Báthory Square Garden":              45,
		"Negocios - Botica Slug & Jiggers":              46,
		"Negocios - Chez Winnie":                        47,
		"Negocios - Danceteria Rolling Hall":            48,
		"Negocios - El Lux":                             49,
		"Negocios - FLEUR":                              50,
		"Negocios - Flourish & Blotts":                  51,
		"Negocios - Luxxuria":                           52,
		"Negocios - Moonlight Shadow Planetary":         53,
		"Negocios - Mystic Momentum":                    54,
		"Negocios - Nym's Treasure":                     55,
		"Negocios - Peonie's Ribbon":                    56,
		"Negocios - Plants & Seeds":                     57,
		"Negocios - Ragnarok":                           58,
		"Negocios - Rose":                               59,
		"Negocios - Royal Vauxhall Tavern":              60,
		"Negocios - Sacred Lotus":                       61,
		"Negocios - Saint Ellis Hospital":               62,
		"Negocios - Sortilegios Weasley":                63,
		"Negocios - Tierra y Cristal":                   64,
		"Negocios - Vinos Zabini":                       65,
		"Negocios - Wizarding World Curse":              66,
		"Logros - Logros de Colaborador del Mes":        67,
		"Logros - Logros de Personaje del Mes":          68,
		"Logros - Logros de Awards":                     69,
		"Logros - Logros de Duelos":                     70,
		"Logros - Logros de Misiones":                   71,
		"Logros - Logros de Expediciones":               72,
		"Logros - Logros de Pociones":                   73,
		"Logros - Logros de Nimbus":                     74,
		"Logros - Logros de Situaciones":                75,
		"Logros - Logros de Cámara de Creación Mágica":  76,
		"Logros - Logros de Rituales":                   77,
		"Logros - Logros de San Mungo":                  78,
		"Logros - Logros de Costello":                   79,
		"Logros - Logros de Colección de Cromos":        80,
		"Logros - Logros de Top Posteadores":            81,
		"Logros - Logros de Empleo":                     82,
		"Logros - Logros de Hall of Fame":               83,
		"Logros - Logros de Torneo de Duelos":           84,
		"Logros - Logros de Torneo de Pociones":         85,
		"Logros - Logros de Torneo de Quidditch Libre":  86,
		"Logros - Logros de Desafío Hogwarts":           87,
		"Logros - Logros de Mortífagos":                 88,
		"Logros - Logros de Aurores":                    89,
	}
	sql := ""
	scvLines := "name|imgurUrl\n"

	for _, cat := range categories {

		for _, item := range cat.Items {
			name := item.Name
			image := "migrated/" + item.Filename
			description := item.Description
			price := item.Price
			stock := "999"
			if price == "" {
				stock = "0"
			}
			module := "0"
			info1 := "0"
			info2 := "0"
			info3 := "0"
			info4 := "0"
			input_needed := "0"
			can_use_item := "0"
			delete_after_use := "0"
			catid := strconv.Itoa(catList[item.Shop+" - "+item.Category])
			status := "0"
			itemlimit := "0"

			//ARREGLO MASCARAS DE MORTIFAGOS
			switch item.Name {
			case "Máscaras 1":
				name = "Máscara de Mortífago 1"
				image = "migrated/Objetos____Mortifagos____Mascaras_1.png"
			case "Máscaras 2":
				name = "Máscara de Mortífago 2"
				image = "migrated/Objetos____Mortifagos____Mascaras_2.png"
			case "Máscaras 3":
				name = "Máscara de Mortífago 3"
				image = "migrated/Objetos____Mortifagos____Mascaras_3.png"
			case "Máscaras 4":
				name = "Máscara de Mortífago 4"
				image = "migrated/Objetos____Mortifagos____Mascaras_4.png"
			case "Máscaras 5":
				name = "Máscara de Mortífago 5"
				image = "migrated/Objetos____Mortifagos____Mascaras_5.png"
			case "Máscaras 6":
				name = "Máscara de Mortífago 6"
				image = "migrated/Objetos____Mortifagos____Mascaras_6.png"
			case "Máscaras 7":
				name = "Máscara de Mortífago 7"
				image = "migrated/Objetos____Mortifagos____Mascaras_7.png"
			case "Máscaras 8":
				name = "Máscara de Mortífago 8"
				image = "migrated/Objetos____Mortifagos____Mascaras_8.png"
			case "Máscaras 9":
				name = "Máscara de Mortífago 9"
				image = "migrated/Objetos____Mortifagos____Mascaras_9.png"
			case "Máscaras 10":
				name = "Máscara de Mortífago 10"
				image = "migrated/Objetos____Mortifagos____Mascaras_10.png"
			}

			query := fmt.Sprintf("INSERT INTO `smf_stshop_items` (`name`, `image`, `description`, `price`, `stock`, `module`, `info1`, `info2`, `info3`, `info4`, `input_needed`, `can_use_item`, `delete_after_use`, `catid`, `status`, `itemlimit`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');", name, image, description, price, stock, module, info1, info2, info3, info4, input_needed, can_use_item, delete_after_use, catid, status, itemlimit)
			sql += query + "\n"
			scvLines += fmt.Sprintf("%s|%s\n", name, item.ImgUrl)
		}

	}
	//writes csv file
	err := os.WriteFile(csvOutputFile, []byte(scvLines), 0644)
	util.Panic(err)

	return sql
}

func generateUrlList(category parser.ShopCategory) error {

	// Crear el nombre del archivo usando category.Name
	fileName := "ImgurURLs.txt"

	// Create or open and append to it
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo: %v", err)
	}
	defer file.Close()

	// Escribir cada URL en el archivo
	for _, item := range category.Items {
		_, err := file.WriteString(item.ImgUrl + "\n")
		if err != nil {
			return fmt.Errorf("no se pudo escribir en el archivo: %v", err)
		}
	}

	//fmt.Printf(category.Name+": Archivo %s generado exitosamente.\n", fileName)
	return nil
}

func DownloadItemPics(category parser.ShopCategory) error {
	for _, item := range category.Items {
		dir := Shop + " - " + category.Name
		time.Sleep(2 * time.Second)
		downloadImage(item.ImgUrl, dir, item.NewImgUrl)
	}
	return nil
}

func downloadImage(url, saveDir, newName string) error {
	// Obtener la extensión del archivo a partir de la URL
	parts := strings.Split(url, "/")
	originalName := parts[len(parts)-1]
	extension := filepath.Ext(originalName)

	// Crear el nombre completo del archivo con la nueva extensión
	imageName := newName + extension

	// Hacer la petición GET a la URL de la imagen
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar que la petición fue exitosa
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s no se pudo descargar, código de estado %d", originalName, resp.StatusCode)
	}

	// Crear el archivo en el directorio de destino
	filePath := filepath.Join(saveDir, imageName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir el contenido de la imagen en el archivo
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("%s descargada exitosamente como %s.\n", originalName, imageName)
	return nil
}

func processImages(category *parser.ShopCategory) error {
	// Directorio donde se encuentran las imágenes descargadas
	picsDir := "./Pics"
	// Directorio de destino para las imágenes copiadas
	destDir := "./" + imageOutputDir

	// Crear el directorio de destino si no existe
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("no se pudo crear el directorio %s: %v", destDir, err)
		}
	}

	// Variables para el reporte
	totalItems := len(category.Items)
	processedCount := 0
	var failedItems []string

	// Iterar sobre cada item en la categoría
	for i := range category.Items {
		item := &category.Items[i]
		// Extraer el nombre del archivo y la extensión del campo ImgUrl
		parts := strings.Split(item.ImgUrl, "/")
		originalFileName := parts[len(parts)-1]
		extension := filepath.Ext(originalFileName)

		if extension == ".gif" {
			extension = ".mp4"
			//remove extension from originalFileName
			originalFileName = strings.TrimSuffix(originalFileName, filepath.Ext(originalFileName))
			originalFileName = originalFileName + extension
		}

		originalFileName = strings.TrimSpace(originalFileName)
		originalFileName = strings.ReplaceAll(originalFileName, "á", "a")
		originalFileName = strings.ReplaceAll(originalFileName, "é", "e")
		originalFileName = strings.ReplaceAll(originalFileName, "í", "i")
		originalFileName = strings.ReplaceAll(originalFileName, "ó", "o")
		originalFileName = strings.ReplaceAll(originalFileName, "ú", "u")
		originalFileName = strings.ReplaceAll(originalFileName, " ", "+")

		// Ruta completa del archivo original en la carpeta PICS
		originalFilePath := filepath.Join(picsDir, originalFileName)

		// Verificar si el archivo existe en PICS
		if _, err := os.Stat(originalFilePath); os.IsNotExist(err) {
			failedItems = append(failedItems, item.Name+" - "+originalFileName)
			continue
		}

		// Crear el nuevo nombre de archivo usando NewImgUrl y manteniendo la extensión original
		catName := item.Category + "__"
		catName = strings.ReplaceAll(catName, " ", "_")
		catName = strings.TrimSpace(catName)
		catName = strings.ReplaceAll(catName, "á", "a")
		catName = strings.ReplaceAll(catName, "é", "e")
		catName = strings.ReplaceAll(catName, "í", "i")
		catName = strings.ReplaceAll(catName, "ó", "o")
		catName = strings.ReplaceAll(catName, "ú", "u")
		catName = strings.ReplaceAll(catName, " ", "+")
		catName = regexp.MustCompile("[^a-zA-Z0-9_+]+").ReplaceAllString(catName, "")

		shopName := item.Shop + "__"
		shopName = strings.ReplaceAll(shopName, " ", "_")
		shopName = strings.TrimSpace(shopName)
		shopName = strings.ReplaceAll(shopName, "á", "a")
		shopName = strings.ReplaceAll(shopName, "é", "e")
		shopName = strings.ReplaceAll(shopName, "í", "i")
		shopName = strings.ReplaceAll(shopName, "ó", "o")
		shopName = strings.ReplaceAll(shopName, "ú", "u")
		shopName = strings.ReplaceAll(shopName, " ", "+")
		shopName = regexp.MustCompile("[^a-zA-Z0-9_+]+").ReplaceAllString(shopName, "")

		itemName := item.Name
		itemName = strings.ReplaceAll(itemName, " ", "_")
		itemName = strings.TrimSpace(itemName)
		itemName = strings.ReplaceAll(itemName, "á", "a")
		itemName = strings.ReplaceAll(itemName, "é", "e")
		itemName = strings.ReplaceAll(itemName, "í", "i")
		itemName = strings.ReplaceAll(itemName, "ó", "o")
		itemName = strings.ReplaceAll(itemName, "ú", "u")
		itemName = strings.ReplaceAll(itemName, " ", "+")
		itemName = regexp.MustCompile("[^a-zA-Z0-9_+]+").ReplaceAllString(itemName, "")

		prefix := shopName + "__" + catName + "__"

		item.NewImgUrl = itemName
		newFileName := prefix + item.NewImgUrl + extension
		newFilePath := filepath.Join(destDir, newFileName)

		switch item.Name {
		case "Máscaras 1":
			item.Name = "Máscara de Mortífago 1"
			newFileName = "Objetos____Mortifagos____Mascaras_1.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_1.png")
		case "Máscaras 2":
			item.Name = "Máscara de Mortífago 2"
			newFileName = "Objetos____Mortifagos____Mascaras_2.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_2.png")
		case "Máscaras 3":
			item.Name = "Máscara de Mortífago 3"
			newFileName = "Objetos____Mortifagos____Mascaras_3.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_3.png")
		case "Máscaras 4":
			item.Name = "Máscara de Mortífago 4"
			newFileName = "Objetos____Mortifagos____Mascaras_4.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_4.png")
		case "Máscaras 5":
			item.Name = "Máscara de Mortífago 5"
			newFileName = "Objetos____Mortifagos____Mascaras_5.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_5.png")
		case "Máscaras 6":
			item.Name = "Máscara de Mortífago 6"
			newFileName = "Objetos____Mortifagos____Mascaras_6.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_6.png")
		case "Máscaras 7":
			item.Name = "Máscara de Mortífago 7"
			newFileName = "Objetos____Mortifagos____Mascaras_7.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_7.png")
		case "Máscaras 8":
			item.Name = "Máscara de Mortífago 8"
			newFileName = "Objetos____Mortifagos____Mascaras_8.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_8.png")
		case "Máscaras 9":
			item.Name = "Máscara de Mortífago 9"
			newFileName = "Objetos____Mortifagos____Mascaras_9.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_9.png")
		case "Máscaras 10":
			item.Name = "Máscara de Mortífago 10"
			newFileName = "Objetos____Mortifagos____Mascaras_10.png"
			originalFilePath = filepath.Join(picsDir, "Objetos____Mortifagos____Mascaras_10.png")
		}

		item.Filename = newFileName

		// Copiar el archivo original a la carpeta de destino
		err := copyFile(originalFilePath, newFilePath)
		if err != nil {
			failedItems = append(failedItems, originalFileName)
			continue
		}

		// Incrementar el contador de imágenes procesadas
		processedCount++
	}

	// Generar el contenido del reporte

	if processedCount != totalItems {
		fmt.Println(fmt.Sprintf(category.Name+" %d de %d imágenes procesadas exitosamente.", processedCount, totalItems))
	}

	if len(failedItems) > 0 {
		fmt.Println("Imágenes que fallaron en ser procesadas:")
		for _, failedItem := range failedItems {
			fmt.Println("\t" + failedItem)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copiar el contenido de sourceFile a destinationFile
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Forzar a que los datos sean escritos al disco
	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
