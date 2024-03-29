package parser

import (
	"localdev/dobby-server/internal/pkg/util"
	"strconv"
	"strings"
)

func ParseDiceRoll(dicerolls []string) []*Dice {
	var dices []*Dice

	for _, diceroll := range dicerolls {
		resultStr := strings.TrimSpace(strings.Split(diceroll, ":")[len(strings.Split(diceroll, ":"))-1])
		//convert result from string to int
		result, err := strconv.Atoi(resultStr)
		util.Panic(err)
		dice := &Dice{
			DiceLine: diceroll,
			Result:   result,
		}
		dices = append(dices, dice)
	}
	return dices
}

func GetTandTopicNameFromViewTopicUrl(viewTopicUrl string) (string, string) {
	// Extracting the value of 't'
	startIndex := strings.Index(viewTopicUrl, "t=") + len("t=")
	endIndex := strings.Index(viewTopicUrl[startIndex:], "&")
	tValue := viewTopicUrl[startIndex : startIndex+endIndex]

	// Extracting the value of 'topic_name'
	topicNameIndex := strings.Index(viewTopicUrl, "&") + 1
	topicNameEndIndex := strings.Index(viewTopicUrl[topicNameIndex:], "#")
	topicName := viewTopicUrl[topicNameIndex : topicNameIndex+topicNameEndIndex]

	return tValue, topicName
}
