package handler

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
	"localdev/dobby-server/internal/pkg/util"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func GetCurrentSessionAndTool(c echo.Context, t *tool.Tool) (*model.UserSession, *tool.Tool) {
	sess, err := session.Get("session", c)
	if sess.IsNew {
		util.LongPrintlnPrintln("New session")
	} else if err != nil {
		util.Panic(err)
	}

	userSession, ok := sess.Values["user_session"].(*model.UserSession)
	if !ok {
		userSession = &model.UserSession{
			IsLoggedIn:    false,
			Username:      nil,
			Initials:      nil,
			LoginDatetime: nil,
		}
	}

	if userSession.ForumCookies != nil && len(userSession.ForumCookies) > 0 && IsCookiesFromUser(userSession) {
		forumUrl, err := url.Parse("https://www.hogwartsrol.com")
		util.Panic(err)
		jar, err := cookiejar.New(nil)
		util.Panic(err)

		client := &http.Client{
			Jar: jar,
		}
		var cookies []*http.Cookie
		for _, cookieEntry := range userSession.ForumCookies {
			cookies = append(cookies, ConvertCookieEntryToCookie(cookieEntry))
		}

		jar.SetCookies(forumUrl, cookies)

		t.PostSecret1 = userSession.PostSecret1
		t.PostSecret2 = userSession.PostSecret2
		t.Client = client
	}

	return userSession, t
}

func IsCookiesFromUser(userSession *model.UserSession) bool {
	isCookieFromUser := false
	if userSession.ForumCookies != nil && len(userSession.ForumCookies) > 0 {
		for _, cookieEntry := range userSession.ForumCookies {
			if cookieEntry.Name == "username" {
				isCookieFromUser = cookieEntry.Value == userSession.User.Username
				util.LongPrintlnPrintln("Username: " + cookieEntry.Value)
			}
		}
	}
	return isCookieFromUser
}

func GetForumCookiesFromClient(client *http.Client, username string) []model.CookieEntry {
	forumUrl, err := url.Parse("https://www.hogwartsrol.com")
	util.Panic(err)
	forumCookies := client.Jar.Cookies(forumUrl)

	var forumCookiesEntries []model.CookieEntry
	for _, cookie := range forumCookies {
		forumCookiesEntries = append(forumCookiesEntries, model.CookieEntry{
			Name:   cookie.Name,
			Value:  cookie.Value,
			Domain: cookie.Domain,
		})
	}
	customCookie := model.CookieEntry{
		Name:   "username",
		Value:  username,
		Domain: "www.hogwartsrol.com",
	}
	forumCookiesEntries = append(forumCookiesEntries, customCookie)

	return forumCookiesEntries
}

func ConvertCookieEntryToCookie(cookieEntry model.CookieEntry) *http.Cookie {
	return &http.Cookie{
		Name:   cookieEntry.Name,
		Value:  cookieEntry.Value,
		Domain: cookieEntry.Domain,
	}
}
