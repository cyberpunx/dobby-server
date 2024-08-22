package main

import (
	"bufio"
	"encoding/csv"
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
	categoriasInput   = "smf_stshop_categories.csv"

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
		//get prices from shop
		categories = parser.ParseShop(shopHtmlStr, Shop, categories)

		for _, cat := range categories {
			err = processImages(&cat)
			util.Panic(err)
		}
	}

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

func readCategoriesFromCsv() map[string]int {
	csvFile, err := os.Open(categoriasInput)
	util.Panic(err)
	defer csvFile.Close()

	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	csvReader.Comma = csvDelimiter
	csvReader.LazyQuotes = true

	categories := make(map[string]int)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		//if first line, skip
		if record[0] == "catid" {
			continue
		}
		util.Panic(err)
		catid := record[0]
		//remove quotes
		catid = strings.ReplaceAll(catid, "\"", "")
		catid = strings.TrimSpace(catid)
		id, err := strconv.Atoi(catid)
		util.Panic(err)

		catName := record[1]
		//remove quotes
		catName = strings.ReplaceAll(catName, "\"", "")
		catName = strings.TrimSpace(catName)

		categories[catName] = id
	}

	return categories
}

func generateSqlQuery(categories []parser.ShopCategory) string {
	catList := readCategoriesFromCsv()
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
			if catid == "0" {
				fmt.Println("Categoria no encontrada: " + item.Shop + " - " + item.Category)
			}
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
