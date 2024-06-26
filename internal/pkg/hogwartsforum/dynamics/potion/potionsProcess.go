package potion

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	gsheet "localdev/dobby-server/internal/pkg/gsheet"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics"
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/util"
	"strconv"
	"strings"
	"time"
)

const (
	Player1   = "Player1"
	Player2   = "Player2"
	Moderator = "Moderator"
	Other     = "Moderator"

	StatusSuccess                Status = "Success"
	StatusFailButMightSucceed    Status = "FailButMightSucceed"
	StatusFail                   Status = "Fail"
	StatusWaitingPlayer1         Status = "WaitingPlayer1"
	StatusWaitingPlayer2         Status = "WaitingPlayer2"
	StatusWaitingPlayer1OnDayOff Status = "WaitingPlayer1OnDayOff"
	StatusWaitingPlayer2OnDayOff Status = "WaitingPlayer2OnDayOff"

	DayOffExtraHours = 24
	ModMaxBonus      = 5

	TemplatePath = "internal/pkg/hogwartsforum/dynamics/templates/"

	FailBecauseOfTime       FailType = "Time"
	FailBecauseOfScore      FailType = "Score"
	FailBecauseOfEditedDice FailType = "EditedDice"
)

type Status string

type Potion struct {
	Name        string
	Ingredients []string
	TargetScore int
	TurnLimit   int
}

type PotionsUser struct {
	Name        string
	House       string
	Role        string
	ProfileUrl  string
	PlayerBonus int
	Posts       []*parser.Post
}
type PotionClubReport struct {
	Player1     PotionsUser
	Player2     PotionsUser
	Moderator   PotionsUser
	Other       []PotionsUser
	Thread      parser.Thread
	Potion      Potion
	Status      Status
	Score       PotionClubScoreBoard
	TurnLimit   int
	TimeLimit   int
	Turns       []PotionClubTurn
	ElapsedTime time.Duration
}

type ModMsgPotionFailData struct {
	RewardedPlayer            string
	RewardedPlayerGoldAmount  int
	RewardedPlayerHousePoints int
	RewardedPlayerHouse       string
	PenalizedPlayer           string
	PenalizedPlayerHouse      string
}

type ModMsgPotionSuccessData struct {
	DiceScoreSum      int
	ModeratorMalus    int
	ModeratorBonus    int
	PlayersTotalBonus int
	TotalScore        int
	TargetScore       int
	Player1           string
	Player2           string
	PotionIcon        template.HTML
	Player1House      string
	Player2House      string
}

type ModMsgNewPotionData struct {
	Player1     string
	Player2     string
	PotionName  string
	TurnLimit   int
	TargetScore int
	Ingredients []string
}

type PotionClubTurn struct {
	Player         PotionsUser
	Number         int
	DiceValue      int
	OnTime         bool
	DayOffUsed     bool
	TurnDatePosted time.Time
	TurnDateLimit  time.Time
	TimeElapsed    time.Duration
}

type FailType string

type PotionClubScoreBoard struct {
	ReportFailed      ModMsgPotionFailData
	ReportSucced      ModMsgPotionSuccessData
	DiceScoreSum      int
	ModeratorBonus    int
	ModeratorMalus    int
	Player1Bonus      int
	Player2Bonus      int
	PlayersTotalBonus int
	TargetScore       int
	TotalScore        int
	Success           bool
	FailureReason     FailType
	ModMessage        string
}

func getPotionFromThread(thread parser.Thread) *Potion {
	potionPostHtml := thread.Posts[0].Content

	reader := strings.NewReader(potionPostHtml)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil
	}

	var name string
	var turnLimit string
	var targetScore string
	var ingredients []string
	var targetScoreInt int
	var turnLimitInt int

	potionInfo := doc.Find("div.xxEDV").Last()
	potionInfo.Find("li").Each(func(i int, liSelection *goquery.Selection) {
		if i == 0 {
			//Potion Name
			name = liSelection.Text()
			name = strings.Split(name, ":")[1]
			name = strings.TrimSpace(name)
		} else if i == 1 {
			//Potion TurnLimit
			turnLimit = liSelection.Text()
			turnLimit := strings.Split(strings.Split(turnLimit, ":")[1], " ")[1]
			turnLimitInt, _ = strconv.Atoi(turnLimit)
		} else if i == 2 {
			//Potion TargetScore
			targetScore = liSelection.Text()
			targetScore := strings.Split(strings.Split(targetScore, ":")[1], " ")[1]
			targetScoreInt, _ = strconv.Atoi(targetScore)
		} else {
			ingredients = append(ingredients, liSelection.Text())
		}
	})

	return &Potion{
		Name:        name,
		Ingredients: ingredients,
		TargetScore: targetScoreInt,
		TurnLimit:   turnLimitInt,
	}
}
func identifyRolesOnThread(thread parser.Thread) (player1 PotionsUser, player2 PotionsUser, moderator PotionsUser, other []PotionsUser) {
	moderator.Name = thread.Author.Username
	moderator.Role = Moderator

	p1name, p1url, p2name, p2url := parser.GetPotionPlayers(thread.Posts[0].Content)
	player1.Name = p1name
	player1.ProfileUrl = p1url
	player1.Role = Player1
	player2.Name = p2name
	player2.ProfileUrl = p2url
	player2.Role = Player2

	for _, post := range thread.Posts {
		if post.Author.Username == p1name && player1.House == "" {
			player1.House = post.Author.House
		}
		if post.Author.Username == p2name && player2.House == "" {
			player2.House = post.Author.House
		}
	}

	for _, post := range thread.Posts {
		if post.Author.Username != moderator.Name && post.Author.Username != player1.Name && post.Author.Username != player2.Name {
			other = append(other, PotionsUser{Name: post.Author.Username, Role: Other})
		}
	}

	for _, post := range thread.Posts {
		if post.Author.Username == player1.Name {
			player1.Posts = append(player1.Posts, post)
		} else if post.Author.Username == player2.Name {
			player2.Posts = append(player2.Posts, post)
		} else if post.Author.Username == moderator.Name {
			moderator.Posts = append(moderator.Posts, post)
		} else {
			for _, otherUser := range other {
				if post.Author.Username == otherUser.Name {
					otherUser.Posts = append(otherUser.Posts, post)
				}
			}
		}
	}

	return
}
func isPlayer(post parser.Post, player1, player2 PotionsUser) bool {
	return post.Author.Username == player1.Name || post.Author.Username == player2.Name
}
func findPlayerAndRole(username string, player1, player2, moderator PotionsUser, others []PotionsUser) (*PotionsUser, string, *PotionsUser) {
	if username == player1.Name {
		return &player1, player1.Role, &player2
	} else if username == player2.Name {
		return &player2, player2.Role, &player1
	} else if username == moderator.Name {
		return &moderator, moderator.Role, nil
	} else {
		for _, otherUser := range others {
			if username == otherUser.Name {
				return &otherUser, otherUser.Role, nil
			}
		}
	}
	return nil, "", nil
}
func removeOtherPostsFromThread(player1 PotionsUser, player2 PotionsUser, moderator PotionsUser, other []PotionsUser, thread parser.Thread) parser.Thread {
	var threadWithoutOthers parser.Thread
	threadWithoutOthers = thread
	threadWithoutOthers.Posts = nil

	for _, post := range thread.Posts {
		if post.Author.Username == player1.Name || post.Author.Username == player2.Name || post.Author.Username == moderator.Name {
			threadWithoutOthers.Posts = append(threadWithoutOthers.Posts, post)
		}
	}

	return threadWithoutOthers
}

func PotionGetReportFromThread(forumDynamic dynamics.ForumDynamic, rawThread parser.Thread, timeLimitHours, turnLimit int, forumDateTime time.Time, daysOff *[]gsheet.DayOff, playerBonus *[]gsheet.PlayerBonus) PotionClubReport {
	timeThreshold := time.Duration(timeLimitHours) * time.Hour
	potion := getPotionFromThread(rawThread)
	player1, player2, moderator, others := identifyRolesOnThread(rawThread)
	threadWithoutOthers := removeOtherPostsFromThread(player1, player2, moderator, others, rawThread)
	playerPostCount := make(map[string]int)
	lastPostTime := *threadWithoutOthers.Created
	turnCount := 1
	postCount := 1
	//postDice := ""
	diceTotal := 0
	postOnTime := false
	threadLastPost := *threadWithoutOthers.Posts[len(threadWithoutOthers.Posts)-1]
	p1Bonus := gsheet.GetPlayerBonusForUsername(playerBonus, player1.Name)
	p2Bonus := gsheet.GetPlayerBonusForUsername(playerBonus, player2.Name)
	player1.PlayerBonus = p1Bonus
	player2.PlayerBonus = p2Bonus

	result := PotionClubReport{
		Player1:   player1,
		Player2:   player2,
		Moderator: moderator,
		Other:     others,
		Thread:    rawThread,
		Potion:    *potion,
		Status:    StatusWaitingPlayer1,
		Score: PotionClubScoreBoard{
			DiceScoreSum:      0,
			ModeratorBonus:    0,
			ModeratorMalus:    0,
			Player1Bonus:      player1.PlayerBonus,
			Player2Bonus:      player2.PlayerBonus,
			PlayersTotalBonus: player1.PlayerBonus + player2.PlayerBonus,
			TargetScore:       potion.TargetScore,
			TotalScore:        0,
			Success:           false,
			ModMessage:        "",
		},
		TurnLimit: turnLimit,
		TimeLimit: timeLimitHours,
		Turns:     nil,
	}

	for _, post := range threadWithoutOthers.Posts {
		postUser := post.Author.Username
		postPlayer, postRole, notPostPlayer := findPlayerAndRole(postUser, player1, player2, moderator, others)

		isPlayerFlag := postPlayer.Name == player1.Name || postPlayer.Name == player2.Name
		dayOffUsed := false

		if isPlayerFlag {
			playerPostCount[postUser]++
			if postRole == Player1 {
				result.Status = StatusWaitingPlayer2
			} else if postRole == Player2 {
				result.Status = StatusWaitingPlayer1
			}

			//if post is edited, the potion fails automatically
			if post.Edited != nil {
				result.Status = StatusFail
				result.Score.Success = false
				result.Score.FailureReason = FailBecauseOfEditedDice
				result.Score.TargetScore = potion.TargetScore
				generatePotionFailedReport(postPlayer.Name, &result)
				result.Score.ModMessage = generateModMessage(forumDynamic, result)
			}

			postOnTime = util.IsDateWithinTimeLimit(*post.Created, lastPostTime, timeThreshold)
			dateLimit := lastPostTime.Add(timeThreshold)
			// if player post is out of time, check if the player used a day off
			if !postOnTime {
				// check for player day off on google sheet and if it's within the time limit
				playerDayOff := gsheet.FindDayOffForPLayerBetweenDates(daysOff, postPlayer.Name, lastPostTime, dateLimit)
				if playerDayOff != nil {
					// check again considering the extra hours of day off
					if util.IsDateWithinTimeLimit(*post.Created, lastPostTime, timeThreshold+time.Hour*DayOffExtraHours) {
						dayOffUsed = true
					}
				}
			}

			postDiceValue := 0
			if len(post.Dices) != 1 {
				//postDice = "N/A"
			} else {
				postDiceValue = post.Dices[0].Result
				//postDice = config.Yellow + strconv.Itoa(post.Dices[0].Result) + config.Reset
				diceTotal += postDiceValue
				result.Score.DiceScoreSum += postDiceValue
			}

			turn := PotionClubTurn{
				Player:         *postPlayer,
				Number:         turnCount,
				DiceValue:      postDiceValue,
				OnTime:         postOnTime,
				DayOffUsed:     dayOffUsed,
				TurnDatePosted: *post.Created,
				TurnDateLimit:  dateLimit,
				TimeElapsed:    post.Created.Sub(lastPostTime),
			}
			result.Turns = append(result.Turns, turn)
			lastPostTime = *post.Created
		}

		if threadLastPost.Id == post.Id && isPlayerFlag {
			elapsedTime := forumDateTime.Sub(*post.Created)
			dateLimit := lastPostTime.Add(timeThreshold)
			if elapsedTime > timeThreshold {
				// check for player day off on google sheet and if it's within the time limit
				playerDayOff := gsheet.FindDayOffForPLayerBetweenDates(daysOff, notPostPlayer.Name, lastPostTime, dateLimit)
				if playerDayOff != nil {
					// check again considering the extra hours of day off
					if elapsedTime < timeThreshold+time.Hour*DayOffExtraHours {
						if notPostPlayer.Name == player1.Name {
							result.Status = StatusWaitingPlayer1OnDayOff
						} else {
							result.Status = StatusWaitingPlayer2OnDayOff
						}
					}
				} else {
					result.Status = StatusFail
					result.Score.Success = false
					result.Score.FailureReason = FailBecauseOfTime
					result.Score.TargetScore = potion.TargetScore
					generatePotionFailedReport(notPostPlayer.Name, &result)
					result.Score.ModMessage = generateModMessage(forumDynamic, result)
				}
			}
		}

		if turnCount == turnLimit && result.Status == StatusWaitingPlayer1 {
			result.Score.DiceScoreSum = diceTotal
			result.Score.TotalScore = diceTotal + result.Score.ModeratorBonus + result.Score.ModeratorMalus + result.Score.Player1Bonus + result.Score.Player2Bonus

			if diceTotal >= potion.TargetScore {
				result.Status = StatusSuccess
				result.Score.Success = true
				result.Score.TargetScore = potion.TargetScore
				generatePotionSuccessReport(&result)
			} else if diceTotal+ModMaxBonus >= potion.TargetScore {
				result.Status = StatusFailButMightSucceed
				result.Score.Success = true
				result.Score.TargetScore = potion.TargetScore
				generatePotionFailedReport(postPlayer.Name, &result)
			} else if diceTotal+ModMaxBonus < potion.TargetScore {
				result.Status = StatusFail
				result.Score.FailureReason = FailBecauseOfScore
				result.Score.Success = false
				result.Score.TargetScore = potion.TargetScore
				generatePotionFailedReport(postPlayer.Name, &result)
			}
			result.Score.ModMessage = generateModMessage(forumDynamic, result)
		}

		postCount++
		if playerPostCount[player1.Name] > 0 && playerPostCount[player2.Name] > 0 {
			playerPostCount[player1.Name] = 0
			playerPostCount[player2.Name] = 0
			turnCount++
		}
	}

	//if at least 1 turn is out of time and day off was not used, the potion is a fail
	if result.Status != StatusSuccess {
		for _, turn := range result.Turns {
			if !turn.OnTime {
				if !turn.DayOffUsed {
					result.Status = StatusFail
					result.Score.FailureReason = FailBecauseOfTime
					result.Score.Success = false
					generatePotionFailedReport(turn.Player.Name, &result)
					result.Score.ModMessage = generateModMessage(forumDynamic, result)
				}
			}
		}
	}

	//calculate the elapsed time
	var elapsedTime time.Duration
	var lastTurn PotionClubTurn
	if result.Turns != nil && len(result.Turns) > 0 {
		lastTurn = result.Turns[len(result.Turns)-1]
		turnTime := lastTurn.TurnDatePosted
		elapsedTime = forumDateTime.Sub(turnTime)
		result.ElapsedTime = elapsedTime
	} else {
		lastPost := result.Thread.Posts[len(result.Thread.Posts)-1]
		postTime := *lastPost.Created
		elapsedTime = forumDateTime.Sub(postTime)
		result.ElapsedTime = elapsedTime
	}
	if elapsedTime > timeThreshold && lastTurn.Player.Role != Player2 && lastTurn.Number == turnLimit {
		result.Score.Success = false
		result.Score.FailureReason = FailBecauseOfTime
		result.Score.TargetScore = potion.TargetScore
		if result.Status == StatusWaitingPlayer1 {
			generatePotionFailedReport(player1.Name, &result)
		} else {
			generatePotionFailedReport(player2.Name, &result)
		}
		result.Status = StatusFail
		result.Score.ModMessage = generateModMessage(forumDynamic, result)
	}

	return result
}

func generatePotionFailedReport(postPlayerName string, report *PotionClubReport) {
	if postPlayerName == report.Player1.Name {
		report.Score.ReportFailed.PenalizedPlayer = report.Player1.Name
		report.Score.ReportFailed.PenalizedPlayerHouse = parser.HouseNameWithColor[report.Player1.House]
		report.Score.ReportFailed.RewardedPlayer = report.Player2.Name
		report.Score.ReportFailed.RewardedPlayerHouse = parser.HouseNameWithColor[report.Player2.House]
	} else {
		report.Score.ReportFailed.PenalizedPlayer = report.Player2.Name
		report.Score.ReportFailed.PenalizedPlayerHouse = parser.HouseNameWithColor[report.Player2.House]
		report.Score.ReportFailed.RewardedPlayer = report.Player1.Name
		report.Score.ReportFailed.RewardedPlayerHouse = parser.HouseNameWithColor[report.Player1.House]
	}
	turnsPlayed := report.Turns[len(report.Turns)-1].Number
	if turnsPlayed <= 3 {
		report.Score.ReportFailed.RewardedPlayerGoldAmount = 200
		report.Score.ReportFailed.RewardedPlayerHousePoints = 100
	} else if turnsPlayed <= 6 {
		report.Score.ReportFailed.RewardedPlayerGoldAmount = 250
		report.Score.ReportFailed.RewardedPlayerHousePoints = 150
	} else if turnsPlayed <= 9 {
		report.Score.ReportFailed.RewardedPlayerGoldAmount = 300
		report.Score.ReportFailed.RewardedPlayerHousePoints = 200
	}
}

func generatePotionSuccessReport(report *PotionClubReport) {
	report.Score.ReportSucced = ModMsgPotionSuccessData{
		DiceScoreSum:      report.Score.DiceScoreSum,
		ModeratorMalus:    report.Score.ModeratorMalus,
		ModeratorBonus:    report.Score.ModeratorBonus,
		PlayersTotalBonus: report.Score.PlayersTotalBonus,
		TotalScore:        report.Score.TotalScore,
		TargetScore:       report.Score.TargetScore,
		Player1:           report.Player1.Name,
		Player2:           report.Player2.Name,
		PotionIcon:        PotionIcons[report.Potion.Name],
		Player1House:      parser.HouseNameWithColor[report.Player1.House],
		Player2House:      parser.HouseNameWithColor[report.Player2.House],
	}
}

func generateModMessage(forumDynamic dynamics.ForumDynamic, r PotionClubReport) string {
	var templateFile string
	var data interface{}
	var templateFolder string
	switch forumDynamic {
	case dynamics.DynamicPotion:
		templateFolder = "potionTemplates/"
	case dynamics.DynamicCreationChamber:
		templateFolder = "creationChamberTemplates/"
	}

	if r.Score.Success {
		templateFile = TemplatePath + templateFolder + "success.html"
		data = r.Score.ReportSucced
	} else {
		data = r.Score.ReportFailed
		switch r.Score.FailureReason {
		case FailBecauseOfTime:
			templateFile = TemplatePath + templateFolder + "failed_time.html"
		case FailBecauseOfScore:
			templateFile = TemplatePath + templateFolder + "failed_score.html"
		case FailBecauseOfEditedDice:
			templateFile = TemplatePath + templateFolder + "failed_edited_dice.html"
		}
	}

	// Parse the selected template
	tmpl, err := template.ParseFiles(templateFile)
	util.Panic(err)

	// Execute the template with the report data
	var out bytes.Buffer
	err = tmpl.Execute(&out, data)
	util.Panic(err)

	return out.String()
}

func GenerateNewPotionMessage(potionName, player1, player2 string, turnLimit, targetScore int) string {
	ingredientsLine := PotionIngredients[potionName]
	ingrentsArray := strings.Split(ingredientsLine, "\n")

	data := ModMsgNewPotionData{
		Player1:     player1,
		Player2:     player2,
		PotionName:  potionName,
		TurnLimit:   turnLimit,
		TargetScore: targetScore,
		Ingredients: ingrentsArray,
	}

	// Parse the selected template
	tmpl, err := template.ParseFiles(TemplatePath + "potionTemplates/new.html")
	util.Panic(err)

	// Execute the template with the report data
	var out bytes.Buffer
	err = tmpl.Execute(&out, data)
	util.Panic(err)

	return out.String()
}
