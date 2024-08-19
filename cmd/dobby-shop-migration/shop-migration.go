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


				switch item.Category {
				case "Generales":
					item.Shop = "Objetos"
				case "Armas Blancas":
					item.Shop = "Objetos"
				case "Criaturas":
					item.Shop = "Objetos"
				case "Estudiantes":
					item.Shop = "Objetos"
				case "Hogwarts":
					item.Shop = "Objetos"
				case "Ingredientes de Rituales":
					item.Shop = "Objetos"
				case "Ministerio":
					item.Shop = "Objetos"
				case "Mortífagos":
					item.Shop = "Objetos"
				case "Pociones":
					item.Shop = "Objetos"
				case "Propiedades":
					item.Shop = "Objetos"
				case "Quidditch":
					item.Shop = "Objetos"
				case "Transporte Mágico":
					item.Shop = "Objetos"
				case "San Mungo":
					item.Shop = "Objetos"
				case "Autorizaciones":
					item.Shop = "Objetos"
				case "Premios Misiones":
					item.Shop = "Objetos"
				case "Premios Expediciones":
					item.Shop = "Objetos"
				case "Premios Nimbus":
					item.Shop = "Objetos"
				case "Premios Situaciones":
					item.Shop = "Objetos"
				case "Mini-tramas":
					item.Shop = "Objetos"
				case "San Duende":
					item.Shop = "Objetos"
				case "Hechizos Básicos":
					item.Shop = "Hechizos"
				case "Hechizos de Sanación":
					item.Shop = "Hechizos"
				case "Hechizos de Ataque":
					item.Shop = "Hechizos"
				case "Hechizos de Defensa":
					item.Shop = "Hechizos"
				case "Hechizos de Ataque y Defensa":
					item.Shop = "Hechizos"
				case "Hechizos Aurores":
					item.Shop = "Hechizos"
				case "Hechizos Mortífagos":
					item.Shop = "Hechizos"
				case "Maleficios":
					item.Shop = "Hechizos"
				case "Habilidades Adquiribles":
					item.Shop = "Habilidades"
				case "Habilidades de Licántropos":
					item.Shop = "Habilidades"
				case "Habilidades de Semigigantes":
					item.Shop = "Habilidades"
				case "Habilidades Sirenas":
					item.Shop = "Habilidades"
				case "Habilidades de Vampiros":
					item.Shop = "Habilidades"
				case "Habilidades de Veela":
					item.Shop = "Habilidades"
				case "Habilidades de Híbridos Innatas":
					item.Shop = "Habilidades"
				case "Habilidades de Humanos":
					item.Shop = "Habilidades"
				case "Arcana High Bar":
					item.Shop = "Costello"
				case "Borgin & Burkes":
					item.Shop = "Costello"
				case "El Nox":
					item.Shop = "Costello"
				case "Mortem Gemma":
					item.Shop = "Costello"
				case "Portafolio":
					item.Shop = "Costello"
				case "Aquí te tengo tu cariñito":
					item.Shop = "Negocios"
				case "Báthory Square Garden":
					item.Shop = "Negocios"
				case "Botica Slug & Jiggers":
					item.Shop = "Negocios"
				case "Chez Winnie":
					item.Shop = "Negocios"
				case "Danceteria Rolling Hall":
					item.Shop = "Negocios"
				case "El Lux":
					item.Shop = "Negocios"
				case "FLEUR":
					item.Shop = "Negocios"
				case "Flourish & Blotts":
					item.Shop = "Negocios"
				case "Luxxuria":
					item.Shop = "Negocios"
				case "Moonlight Shadow Planetary":
					item.Shop = "Negocios"
				case "Mystic Momentum":
					item.Shop = "Negocios"
				case "Nym's Treasure":
					item.Shop = "Negocios"
				case "Peonie's Ribbon":
					item.Shop = "Negocios"
				case "Plants & Seeds":
					item.Shop = "Negocios"
				case "Ragnarok":
					item.Shop = "Negocios"
				case "Rose":
					item.Shop = "Negocios"
				case "Royal Vauxhall Tavern":
					item.Shop = "Negocios"
				case "Sacred Lotus":
					item.Shop = "Negocios"
				case "Saint Ellis Hospital":
					item.Shop = "Negocios"
				case "Sortilegios Weasley":
					item.Shop = "Negocios"
				case "Tierra y Cristal":
					item.Shop = "Negocios"
				case "Vinos Zabini":
					item.Shop = "Negocios"
				case "Wizarding World Curse":
					item.Shop = "Negocios"
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
		"Habilidades - Habilidades Adquiribles":         29,
		"Habilidades - Habilidades de Licántropos":      30,
		"Habilidades - Habilidades de Semigigantes":     31,
		"Habilidades - Habilidades Sirenas":             32,
		"Habilidades - Habilidades de Vampiros":         33,
		"Habilidades - Habilidades de Veela":            34,
		"Habilidades - Habilidades de Híbridos Innatas": 35,
		"Habilidades - Habilidades de Humanos":          36,
		"Costello - Arcana High Bar":                    37,
		"Costello - Borgin & Burkes":                    38,
		"Costello - El Nox":                             39,
		"Costello - Mortem Gemma":                       40,
		"Costello - Portafolio":                         41,
		"Negocios - Aquí te tengo tu cariñito":          42,
		"Negocios - Báthory Square Garden":              43,
		"Negocios - Botica Slug & Jiggers":              44,
		"Negocios - Chez Winnie":                        45,
		"Negocios - Danceteria Rolling Hall":            46,
		"Negocios - El Lux":                             47,
		"Negocios - FLEUR":                              48,
		"Negocios - Flourish & Blotts":                  49,
		"Negocios - Luxxuria":                           50,
		"Negocios - Moonlight Shadow Planetary":         51,
		"Negocios - Mystic Momentum":                    52,
		"Negocios - Nym's Treasure":                     53,
		"Negocios - Peonie's Ribbon":                    54,
		"Negocios - Plants & Seeds":                     55,
		"Negocios - Ragnarok":                           56,
		"Negocios - Rose":                               57,
		"Negocios - Royal Vauxhall Tavern":              58,
		"Negocios - Sacred Lotus":                       59,
		"Negocios - Saint Ellis Hospital":               60,
		"Negocios - Sortilegios Weasley":                61,
		"Negocios - Tierra y Cristal":                   62,
		"Negocios - Vinos Zabini":                       63,
		"Negocios - Wizarding World Curse":              64,
	}
	sql := ""
	scvLines := "name|imgurUrl\n"

	for _, cat := range categories {

		for _, item := range cat.Items {
			name := item.Name
			image := imageOutputDir + "/" + item.Filename
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
			case "Máscaras 2":
				name = "Máscara de Mortífago 2"
			case "Máscaras 3":
				name = "Máscara de Mortífago 3"
			case "Máscaras 4":
				name = "Máscara de Mortífago 4"
			case "Máscaras 5":
				name = "Máscara de Mortífago 5"
			case "Máscaras 6":
				name = "Máscara de Mortífago 6"
			case "Máscaras 7":
				name = "Máscara de Mortífago 7"
			case "Máscaras 8":
				name = "Máscara de Mortífago 8"
			case "Máscaras 9":
				name = "Máscara de Mortífago 9"
			case "Máscaras 10":
				name = "Máscara de Mortífago 10"
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
