package util

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	mylogger "localdev/dobby-server/internal/pkg/logger"
	"net"
	"strconv"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func LoadConfig(path string, config interface{}) {
	abs, err := filepath.Abs(path)
	Panic(err)
	b, err := ioutil.ReadFile(abs)
	Panic(err)
	err = json.Unmarshal(b, config)
	Panic(err)
}

func PStr(s string) *string {
	return &s
}

func PStrf(s string, values ...interface{}) *string {
	return PStr(fmt.Sprintf(s, values...))
}

func PInt64(i int64) *int64 {
	return &i
}

func PInt(i int) *int {
	return &i
}

func PFlt32(f float32) *float32 {
	return &f
}

func PFlt64(f float64) *float64 {
	return &f
}

func PTime(t time.Time) *time.Time {
	return &t
}

func PBool(b bool) *bool {
	return &b
}

func SplitDateAndTime(dateTime string) (string, string) {
	parts := strings.Split(dateTime, ", ")
	if len(parts) != 2 {
		return "", ""
	}
	date := parts[0]
	time := parts[1]
	return date, time
}

func GetTimeFromTimeZone(timezone string) (time.Time, error) {
	url := fmt.Sprintf("http://worldtimeapi.org/api/timezone/%s", timezone)

	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	var response struct {
		Datetime string `json:"datetime"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return time.Time{}, err
	}

	inputLayout := "2006-01-02T15:04:05.999999-07:00"
	outputLayout := "01/02/2006 15:04"
	dateTimeParcial, _ := time.Parse(inputLayout, response.Datetime)
	formattedDateTime := dateTimeParcial.Format(outputLayout)

	dateTime, _ := time.Parse(outputLayout, formattedDateTime)
	if err != nil {
		return time.Time{}, err
	}

	return dateTime, nil
}

func AdjustDateTime(currentDateTime time.Time, dateString string) time.Time {
	if dateString == "Hoy" {
		return currentDateTime
	} else if dateString == "Ayer" {
		return currentDateTime.AddDate(0, 0, -1)
	}
	return currentDateTime
}

func AdjustDateTimeToStr(currentDate time.Time, dateString string) string {
	if dateString == "Hoy" {
		return currentDate.Format("02/01/2006")
	} else if dateString == "Ayer" {
		return currentDate.AddDate(0, 0, -1).Format("02/01/2006")
	}
	return dateString
}

func IsDateInBetween(date time.Time, startDate time.Time, endDate time.Time) bool {
	return date.After(startDate) && date.Before(endDate)
}

func IsDateInCurrentMonth(date time.Time) bool {
	currentDate := time.Now()
	return date.Month() == currentDate.Month() && date.Year() == currentDate.Year()
}

func MarshalJsonPretty(i interface{}) []byte {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "    ")
	Panic(encoder.Encode(i))
	return buffer.Bytes()
}

func TrimAndToLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func IsDateWithinTimeLimit(currentTime, lastPostTime time.Time, timeThreshold time.Duration) bool {
	// Check if the current exceeds the time threshold
	if lastPostTime.Add(timeThreshold).Before(currentTime) {
		return false
	} else {
		return true
	}
}

func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func GetInitials(name string) string {
	words := strings.Fields(name)
	var initials string

	if len(words) == 1 {
		initials = string(words[0][0]) + string(words[0][1])
	} else {
		initials = string(words[0][0]) + string(words[1][0])
	}

	return strings.ToUpper(initials)
}

func LongPrintlnPrintln(a ...any) {
	// Convert all arguments into a string slice
	stringArgs := make([]string, len(a))
	for i, arg := range a {
		stringArgs[i] = fmt.Sprint(arg)
	}

	// Join all arguments into a single string
	fullString := strings.Join(stringArgs, " ")

	// Print to stdout using util.LongPrintlnPrintln
	//_, err := fmt.Println(fullString)
	//Panic(err)

	// slogecho logger
	mylogger.GetLogger().Info(fullString)

	// Append to a log file
	file, err := os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Panic(err)
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(fullString)
}

func PrintResponseStatus(status string) {
	LongPrintlnPrintln("Response Status: " + status)
}

type P map[string]interface{}

func Fprint(format string, p P) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

func SaveJsonFile(fileName string, data []byte) error {
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadJsonFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseDate(dateStr string) (time.Time, error) {
	var layouts = []string{
		"Mon 2 Jan - 15:04",       // Miér 10 Abr - 7:53
		"Mon 2 Jan 2006 - 15:04",  // Miér 10 Abr 2024 - 7:53
		"Mon 2 Jan - 15:04:05",    // Miér 10 Abr - 7:53:02
		"02.01.06 15:04",          // 10.04.24 7:53
		"02/01/06, 3:04 pm",       // 10/04/24, 07:53 am
		"Mon 2 Jan 2006, 3:04 pm", // Miér 10 Abr 2024, 7:53 am
		"Mon 2 Jan 2006, 15:04",   // Miér 10 Abr 2024, 07:53
		"Mon Jan 2, 2006 3:04 pm", // Miér Abr 10, 2024 7:53 am
		"Mon Jan 2 2006, 15:04",   // Miér Abr 10 2024, 07:53
		"2nd Jan 2006, 3:04 pm",   // Ajustar según el sufijo necesario
		"2nd Jan 2006, 15:04",     // Ajustar según el sufijo necesario
		"Jan 2nd 2006, 3:04 pm",   // Ajustar según el sufijo necesario
		"Jan 2nd 2006, 15:04",     // Ajustar según el sufijo necesario
		"2/1/2006, 3:04 pm",       // 10/4/2024, 7:53 am
		"2/1/2006, 15:04",         // 10/4/2024, 07:53
		"1/2/2006, 3:04 pm",       // 4/10/2024, 7:53 am
		"1/2/2006, 15:04",         // 4/10/2024, 07:53
		"2006-01-02, 3:04 pm",     // 2024-04-10, 7:53 am
		"2006-01-02, 15:04",       // 2024-04-10, 07:53
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", dateStr)
}

func IsUserDateFormatCorrect(userDateFormat string, forumDateTime time.Time) bool {
	var timeStr, dateStr string
	var parts []string // Para almacenar partes divididas de la cadena.

	datetimeStr := userDateFormat
	if strings.Contains(datetimeStr, "Hoy a las") || strings.Contains(datetimeStr, "Ayer a las") {
		parts = strings.Split(datetimeStr, " ")
		if len(parts) < 4 { // Verificar que haya suficientes partes.
			return false
		}
		timeStr = parts[3]                                     // Tiempo extraído correctamente.
		dateStr = AdjustDateTimeToStr(forumDateTime, parts[0]) // Ajuste de fecha según "Hoy" o "Ayer".
	} else {
		parts = strings.Split(datetimeStr, ",")
		if len(parts) < 2 { // Verificar que haya suficientes partes.
			return false
		}
		dateStr = parts[0]
		timeStr = strings.TrimSpace(parts[1])
	}

	layout := "2/1/2006 15:04"
	_, err := time.Parse(layout, dateStr+" "+timeStr)
	return err == nil
}

func GenerateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetElapsedTime(elapsedTime time.Duration) string {
	hours := int(elapsedTime.Hours())
	minutes := int(elapsedTime.Minutes()) - int(elapsedTime.Hours())*60

	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func IsPortInUse(port int) bool {
	conn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return true
	}
	_ = conn.Close()
	return false
}
