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
	"strings"
	"time"
)

const (
	loginUsername = "Desarrollo"
	loginPassword = "programación2055"
	inputCsvFile  = "testUser.csv"
	csvDelimiter  = ','

	costelloInput     = "Costello"
	habilidadesInput  = "Habilidades"
	hechizosInput     = "Hechizos"
	negociosInput     = "Negocios"
	objetosInput      = "Objetos"
	descricionesInput = "descripciones.html"
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
	fmt.Println(fmt.Sprintf("%s\n", util.MarshalJsonPretty(categories)))

	//array with inputs
	inputs := []string{costelloInput, habilidadesInput, hechizosInput, negociosInput, objetosInput}
	for _, input := range inputs {
		Shop = input

		shopContents, err := os.ReadFile(Shop + ".html")
		util.Panic(err)
		shopHtmlStr := string(shopContents)

		shops := parser.ParseShop(shopHtmlStr, Shop, categories)
		//fmt.Println(fmt.Sprintf("%s\n", util.MarshalJsonPretty(shops)))
		for _, cat := range categories {
			err = generateUrlList(cat)
			util.Panic(err)
		}

		for _, shop := range shops {
			err = processImages(&shop)
			util.Panic(err)
		}

		//fmt.Println(fmt.Sprintf("%s\n", util.MarshalJsonPretty(shops)))
	}

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

	fmt.Printf(category.Name+": Archivo %s generado exitosamente.\n", fileName)
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
	destDir := "./RenamedPics"

	prefix := category.Name + "__"
	prefix = strings.ReplaceAll(prefix, " ", "_")
	prefix = strings.TrimSpace(prefix)
	prefix = strings.ReplaceAll(prefix, "á", "a")
	prefix = strings.ReplaceAll(prefix, "é", "e")
	prefix = strings.ReplaceAll(prefix, "í", "i")
	prefix = strings.ReplaceAll(prefix, "ó", "o")
	prefix = strings.ReplaceAll(prefix, "ú", "u")
	prefix = strings.ReplaceAll(prefix, " ", "+")
	prefix = regexp.MustCompile("[^a-zA-Z0-9_+]+").ReplaceAllString(prefix, "")

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
			failedItems = append(failedItems, originalFileName)
			continue
		}

		// Crear el nuevo nombre de archivo usando NewImgUrl y manteniendo la extensión original
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

	// Abrir el archivo result.txt en modo append o crearlo si no existe
	reportFile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("no se pudo abrir o crear el archivo de reporte: %v", err)
	}
	defer reportFile.Close()

	// Generar el contenido del reporte
	reportContent := fmt.Sprintf(category.Name+" reporte final: %d de %d imágenes procesadas exitosamente.\n", processedCount, totalItems)
	if len(failedItems) > 0 {
		reportContent += "Imágenes que fallaron en ser procesadas:\n"
		for _, failedItem := range failedItems {
			reportContent += failedItem + "\n"
		}
	}

	reportContent += "------------------------\n"

	// Escribir el reporte en el archivo
	_, err = reportFile.WriteString(reportContent)
	if err != nil {
		return fmt.Errorf("no se pudo escribir en el archivo de reporte: %v", err)
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
