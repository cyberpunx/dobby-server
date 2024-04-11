package tool

import (
	"fmt"
	gsheet2 "localdev/dobby-server/internal/pkg/gsheet"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics/chronology"
	"localdev/dobby-server/internal/pkg/hogwartsforum/dynamics/potion"
	parser "localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/util"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	ForumRulesThread = "t24-001-normas-del-foro"
)

func (o *Tool) parseSubforum(subHtml string) []*parser.Thread {
	threadList := parser.GetSubforumThreads(subHtml)

	var threads []*parser.Thread
	for _, thread := range threadList {
		threadUrl := parser.SubGetThreadUrl(thread)
		threadHtml := o.GetThread(threadUrl)
		thread := o.ParseThread(threadHtml)
		threads = append(threads, thread)
	}
	return threads
}

func (o *Tool) ParseThread(threadHtml string) *parser.Thread {
	threadTitle, threadUrl, err := parser.ThreadExtractTitleAndURL(threadHtml)
	util.Panic(err)

	var posts []*parser.Post
	var pagesUrl []string
	pagesUrl = append(pagesUrl, threadUrl)
	for {
		// Extract the post list from the current page
		postList := parser.ThreadListPosts(threadHtml)
		for _, post := range postList {
			post := o.parsePost(post)
			posts = append(posts, post)
		}

		// Check if there is a "next" link in the pagination
		nextPageURL, hasMore := parser.ThreadNextPageURL(threadHtml)

		if !hasMore {
			break // No more pages to fetch
		}

		// Fetch the next page and update the threadHtml
		pagesUrl = append(pagesUrl, nextPageURL)
		nextPageHTML := o.GetThread(nextPageURL)
		threadHtml = nextPageHTML
	}

	if posts == nil || len(posts) == 0 {
		return nil
	}

	firstPostId := posts[0].Id
	var filteredPosts []*parser.Post
	filteredPosts = append(filteredPosts, posts[0])
	for _, post := range posts {
		if post.Id != firstPostId {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return &parser.Thread{
		Title:    threadTitle,
		Url:      threadUrl,
		Author:   posts[0].Author,
		Created:  posts[0].Created,
		LastPost: posts[len(posts)-1],
		Pages:    pagesUrl,
		Posts:    filteredPosts,
	}
}

func (o *Tool) parsePost(postHtml string) *parser.Post {
	postUser := parser.PostGetUserName(postHtml)
	postUserUrl := parser.PostGetUserUrl(postHtml)
	postUserHouse := parser.PostGetUserHouse(postHtml)
	postDateTime := parser.PostGetDateAndTime(postHtml, o.ForumDateTime)
	postEditedDateTime := parser.PostGetEditedDateAndTime(postHtml)
	postUrl := parser.PostGetUrl(postHtml)
	postContent := parser.PostGetContent(postHtml)
	dices := parser.ParseDiceRoll(parser.PostGetDices(postHtml))

	return &parser.Post{
		Url:     postUrl,
		Author:  &parser.User{Username: postUser, Url: postUserUrl, House: postUserHouse},
		Created: postDateTime,
		Edited:  postEditedDateTime,
		Content: postContent,
		Dices:   dices,
		Id:      postUrl[strings.LastIndex(postUrl, "#")+1:],
	}
}

func (o *Tool) processPotionsSubforum(forumDynamic dynamics.ForumDynamic, subforumThreads []*parser.Thread, timeLimit, turnLimit int) []potion.PotionClubReport {
	var reportList []potion.PotionClubReport
	for threadIndex, thread := range subforumThreads {
		util.LongPrintlnPrintln("Processing Thread: [" + strconv.Itoa(threadIndex+1) + "/" + strconv.Itoa(len(subforumThreads)) + "] " + "Thread: " + thread.Title + " (Time: " + strconv.Itoa(timeLimit) + "| Turn: " + strconv.Itoa(turnLimit) + ")")
		report := o.processPotionsThread(forumDynamic, *thread, timeLimit, turnLimit)
		reportList = append(reportList, report)
		//reportJson := util.MarshalJsonPretty(report)
		//util.LongPrintlnPrintln(fmt.Sprintf("%s\n", reportJson))
	}
	return reportList
}

func (o *Tool) processPotionsThread(forumDynamic dynamics.ForumDynamic, thread parser.Thread, timeLimit, turnLimit int) potion.PotionClubReport {
	var report potion.PotionClubReport
	daysOff := o.getGoogleSheetPotionsDayOff()
	playerBonus := o.getGoogleSheetPotionsBonus()
	report = potion.PotionGetReportFromThread(forumDynamic, thread, timeLimit, turnLimit, o.ForumDateTime, daysOff, playerBonus)
	return report
}

func (o *Tool) processChronoMainThread(chronoMainThread parser.Thread, hrTool *Tool) {
	util.LongPrintlnPrintln("=== Chronology Thread Begin ===")
	util.LongPrintlnPrintln("Thread: " + chronoMainThread.Title)

	var chronoLinks []string
	for _, post := range chronoMainThread.Posts {
		chronoLink := parser.PostGetLinks(post.Content)
		chronoLinks = append(chronoLinks, chronoLink...)
	}

	re, err := regexp.Compile(`p\d+`)
	util.Panic(err)
	var cleanedURLs []string
	for _, link := range chronoLinks {
		parsedURL, err := url.Parse(link)
		util.Panic(err)
		parsedURL.Fragment = ""
		urlWithoutFragment := parsedURL.String()
		cleanedUrl := re.ReplaceAllString(urlWithoutFragment, "")
		cleanedURLs = append(cleanedURLs, cleanedUrl)
	}

	var threadListHtml []string
	for _, link := range cleanedURLs {
		chronoThreadtHtml := hrTool.GetThread(link)
		if parser.IsThreadVisible(chronoThreadtHtml) {
			threadListHtml = append(threadListHtml, chronoThreadtHtml)
		}
	}

	var chronoThreads []*chronology.ChronoThread
	for _, threadHtml := range threadListHtml {
		thread := hrTool.ParseThread(threadHtml)
		chronoThread := chronology.ChronoThreadProcessor(*thread)
		chronoThreads = append(chronoThreads, chronoThread)
	}

	chronoReport := chronology.ChronoReport{
		ChronoThreads: chronoThreads,
	}

	stringContents := fmt.Sprintf("%s\n", util.MarshalJsonPretty(chronoReport))
	filename := "output.json"

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write the content to file
	_, err = file.WriteString(stringContents)
	if err != nil {
		log.Fatal(err)
	}

	util.LongPrintlnPrintln("\n")
	util.LongPrintlnPrintln("=== Chronology Thread End === \n")
}

func (o *Tool) getGoogleSheetPotionsDayOff() *[]gsheet2.DayOff {
	rows, err := gsheet2.ReadSheetData(o.SheetService, o.Config.GSheetModeracionId, gsheet2.SheetRangeDaysOff)
	util.Panic(err)
	daysOff := gsheet2.ParseDayOff(rows)
	return &daysOff
}

func (o *Tool) getGoogleSheetPotionsBonus() *[]gsheet2.PlayerBonus {
	rows, err := gsheet2.ReadSheetData(o.SheetService, o.Config.GSheetModeracionId, gsheet2.SheetRangePlayerBonus)
	util.Panic(err)
	playerBonus := gsheet2.ParsePlayerBonus(rows)
	return &playerBonus
}

func (o *Tool) ProcessPotionsSubforumList(forumDynamic dynamics.ForumDynamic, subForumUrls *[]string, timeLimit, turnLimit *int) []potion.PotionClubReport {
	util.LongPrintlnPrintln("Dynamic: " + strings.ToUpper(string(forumDynamic)) + " / TimeLimit: " + strconv.Itoa(*timeLimit) + " / TurnLimit: " + strconv.Itoa(*turnLimit))
	if len(*subForumUrls) == 0 {
		util.LongPrintlnPrintln("No subforums URLs to process")
	}
	var reportMainList []potion.PotionClubReport
	for _, subforumUrl := range *subForumUrls {
		subforumHtml := o.getSubforum(subforumUrl)
		subforumThreads := o.parseSubforum(subforumHtml)
		reportList := o.processPotionsSubforum(forumDynamic, subforumThreads, *timeLimit, *turnLimit)
		reportMainList = append(reportMainList, reportList...)
	}

	return reportMainList
}

func (o *Tool) ProcessPotionsThreadList(forumDynamic dynamics.ForumDynamic, threadsUrls *[]string, timeLimit, turnLimit *int) []potion.PotionClubReport {
	util.LongPrintlnPrintln("\n\n ========= THREADS DE POCIONES =========\n\n")
	if len(*threadsUrls) == 0 {
		util.LongPrintlnPrintln("No Threads URLs to process")
	}
	var reportMainList []potion.PotionClubReport
	for _, threadUrl := range *threadsUrls {
		potionThreadHtml := o.GetThread(threadUrl)
		potionThread := o.ParseThread(potionThreadHtml)
		report := o.processPotionsThread(forumDynamic, *potionThread, *timeLimit, *turnLimit)
		reportMainList = append(reportMainList, report)
	}

	return reportMainList
}

func (o *Tool) GetUserDateTimeFormat() string {
	threadUrl := o.Config.BaseUrl + ForumRulesThread
	threadHtml := o.GetThread(threadUrl)
	dateTime := parser.PostGetDateTime(threadHtml)
	return dateTime
}
