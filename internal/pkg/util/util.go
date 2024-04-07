package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	conf "localdev/dobby-server/internal/app/dobby-server/config"

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

	// Append to a log file
	file, err := os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Panic(err)
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(fullString)
}

func PrintResponseStatus(status string) {
	statusColor := ""
	statusEmoji := ""
	if status == "200 OK" {
		statusColor = conf.Green
		statusEmoji = " " + conf.CheckEmoji + " "
	} else {
		statusColor = conf.Red
		statusEmoji = " " + conf.CrossEmoji + " "
	}
	LongPrintlnPrintln("Response Status: " + statusColor + statusEmoji + " " + status + conf.Reset)
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
