package tool

import (
	"fmt"
	"io/ioutil"
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/util"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func LoginAndGetCookies(user, pass string) (*http.Client, *LoginResponse) {
	util.LongPrintlnPrintln("Logging in with UserSession: " + user)
	params := url.Values{}
	params.Add("username", user)
	params.Add("password", pass)
	params.Add("autologin", `on`)
	params.Add("redirect", ``)
	params.Add("query", ``)
	params.Add("login", `Conectarse`)
	bodyRequest := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://www.hogwartsrol.com/login", bodyRequest)
	util.Panic(err)
	req.Header.Set("Authority", "www.hogwartsrol.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "es")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "_gid=GA1.2.203551841.1698500337; toolbar_state=fa_show; _pbjs_userid_consent_data=6683316680106290; _pubcid=19fb297d-d17d-4729-8dba-f5dfc67ec7bb; trc_cookie_storage=taboola%2520global%253Auser-id%3Ddc9b0bdc-e3d7-4fb3-962b-c8b849fd38b6-tuctc369474; cto_bidid=2LV_hl9JMTY2SERCbUlDUFRSWHd4QnFNYnU2eFFRaEZ1bzMzcVJwbW9Nb1hDOTNURFFjdThycThLMXYzbVBDa0N4YmRza0p0cTNMVm81a2J6eWp1em5EWSUyRlBnJTNEJTNE; cto_bundle=X7xxSl9tM0VvSUFRV1d4alJqOW5NNDNCSmtaJTJGS1d5WW1jbUJwTVVOSDlOcTI5Nk1wNkI4aWJSRnR2NGpueWJaRWNFUnJua0ZwYkElMkJmdUx1bkwybmJ2Ynl4OFJTaXlmbjZMZWxsTkRScGVCTzBkZzJMT2ZJS3NiVXdyNTk0aGRSN1JVbnI; _fa-screen=%7B%22w%22%3A1681%2C%22h%22%3A1058%7D; _gat_gtag_UA_144386270_1=1; _ga_TTF1KWE3G4=GS1.1.1698500337.1.1.1698500422.59.0.0; _ga=GA1.1.1824435064.1698500337")
	req.Header.Set("Origin", "https://www.hogwartsrol.com")
	req.Header.Set("Referer", "https://www.hogwartsrol.com/login?")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	jar, err := cookiejar.New(nil)
	if err != nil {
		// error handling
	}

	client := &http.Client{
		Jar: jar,
	}
	resp, err := client.Do(req)
	util.Panic(err)

	defer resp.Body.Close()
	util.PrintResponseStatus(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// error handling
	}

	var loginResponse LoginResponse
	isLoginCorrect, msg := parser.IsLoginCorrect(string(body))
	if !isLoginCorrect {
		util.LongPrintlnPrintln("ERROR: Usuario y/o Contraseña Incorrectos")
		loginResponse = LoginResponse{
			Success:  util.PBool(false),
			Messaage: &msg,
			Username: nil,
			Initials: nil,
		}
		return nil, &loginResponse
	}
	username := parser.GetUsername(string(body))
	util.LongPrintlnPrintln("Bienvenido: " + username)
	initials := util.GetInitials(username)
	timestamp := time.Now()
	loginResponse = LoginResponse{
		Success:       util.PBool(true),
		Messaage:      &msg,
		Username:      &username,
		Initials:      &initials,
		LoginDatetime: &timestamp,
	}

	return client, &loginResponse
}

func (o *Tool) PostNewThread(subforumId, subject, message string, notify, attachSig bool, hasDice bool) (*parser.Thread, error) {
	util.LongPrintlnPrintln("Posting New Topic: " + subject + " on subforum " + subforumId)

	data := url.Values{
		"subject": {subject},
		"message": {message},
		"post":    {"Enviar"},
		"f":       {subforumId},
		"lt":      {"0"},
		"mode":    {"newtopic"},
		"auth[]":  {*o.PostSecret1},
	}
	data.Add("auth[]", *o.PostSecret2)

	if notify {
		data.Add("notify", "on")
	}

	if attachSig {
		data.Add("attach_sig", "on")
	}

	if hasDice {
		data.Add("post_dice_0", "")
	}

	queryValues := url.Values{}
	queryValues.Add("f", subforumId)
	queryValues.Add("mode", "newtopic")

	baseDomain := o.Config.BaseUrl

	fullUrl := baseDomain + "post?" + queryValues.Encode()
	req, err := http.NewRequest(http.MethodPost, fullUrl, strings.NewReader(data.Encode()))
	util.Panic(err)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.Client.Do(req)
	util.Panic(err)
	defer resp.Body.Close()
	util.PrintResponseStatus(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	util.Panic(err)

	success, viewTopicUrl := parser.IsPostSuccessful(string(body))
	if !success {
		util.LongPrintlnPrintln("ERROR: Could not post new topic")
		return nil, fmt.Errorf("Could not post new topic check logs for more info: \n" + string(body))
	} else {
		util.LongPrintlnPrintln("OK: New Topic Posted: " + viewTopicUrl)
		t, topic_name := parser.GetTandTopicNameFromViewTopicUrl(viewTopicUrl)
		threadBody := o.getThreadByViewTopic(t, topic_name)
		thread := o.ParseThread(threadBody)
		util.LongPrintlnPrintln("New Topic URL: " + thread.Url)
		return thread, nil
	}
}

func (o *Tool) ReplyThread(threadUrl, message string, notify, attachSig bool) (*parser.Thread, error) {

	threadBody := o.GetThread(threadUrl)
	threadTitle, _, err := parser.ThreadExtractTitleAndURL(threadBody)
	util.LongPrintlnPrintln("Replying on topic: " + threadTitle)
	tid, t, lt, auth1, auth2, err := parser.ThreadExtactReplyData(threadBody)

	attachSigStr := "1"
	if !attachSig {
		attachSigStr = "0"
	}

	notifyStr := "1"
	if !notify {
		notifyStr = "0"
	}

	data := url.Values{
		"attach_sig": {attachSigStr},
		"notify":     {notifyStr},
		"message":    {message},
		"post":       {"Enviar"},
		"t":          {t},
		"lt":         {lt},
		"tid":        {tid},
		"mode":       {"reply"},
		"auth[]":     {auth1},
	}
	data.Add("auth[]", auth2)

	queryValues := url.Values{}
	queryValues.Add("t", t)
	queryValues.Add("mode", "reply")

	baseDomain := o.Config.BaseUrl

	fullUrl := baseDomain + "post?" + queryValues.Encode()
	req, err := http.NewRequest(http.MethodPost, fullUrl, strings.NewReader(data.Encode()))
	util.Panic(err)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.Client.Do(req)
	util.Panic(err)
	defer resp.Body.Close()
	util.PrintResponseStatus(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	util.Panic(err)

	success, viewTopicUrl := parser.IsPostSuccessful(string(body))
	if !success {
		util.LongPrintlnPrintln("ERROR: Could not reply topic")
		return nil, fmt.Errorf("Could not reply topic")
	} else {
		util.LongPrintlnPrintln("OK: New Reply Posted: " + viewTopicUrl)
		t, topic_name := parser.GetTandTopicNameFromViewTopicUrl(viewTopicUrl)
		threadBody := o.getThreadByViewTopic(t, topic_name)
		thread := o.ParseThread(threadBody)
		lastPost := thread.Posts[len(thread.Posts)-1]
		util.LongPrintlnPrintln("New Topic URL: " + lastPost.Url)
		return thread, nil
	}
}

func (o *Tool) getSubforum(subUrl string) string {
	util.LongPrintlnPrintln("Getting Sub: " + subUrl)

	baseDomain := o.Config.BaseUrl
	req, err := http.NewRequest("GET", baseDomain+subUrl, nil)
	util.Panic(err)

	resp, err := o.Client.Do(req)
	util.Panic(err)
	defer resp.Body.Close()
	util.PrintResponseStatus(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	util.Panic(err)

	return string(body)
}

func (o *Tool) getForumHome() string {
	util.LongPrintlnPrintln("Getting Home (Get Forum LoginDatetime): ")

	baseDomain := o.Config.BaseUrl
	req, err := http.NewRequest("GET", baseDomain, nil)
	util.Panic(err)

	resp, err := o.Client.Do(req)
	util.Panic(err)
	defer resp.Body.Close()
	util.PrintResponseStatus(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	util.Panic(err)

	return string(body)
}

func (o *Tool) GetThread(threadUrl string) string {
	util.LongPrintlnPrintln("Getting Thread: " + threadUrl)

	baseDomain := o.Config.BaseUrl

	_, err := url.ParseRequestURI(threadUrl)
	if err != nil || !strings.HasPrefix(threadUrl, baseDomain) {
		threadUrl = baseDomain + threadUrl
	}

	req, err := http.NewRequest("GET", threadUrl, nil)
	util.Panic(err)

	resp, err := o.Client.Do(req)
	util.Panic(err)
	defer resp.Body.Close()
	util.PrintResponseStatus(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	util.Panic(err)

	return string(body)
}

func (o *Tool) getThreadByViewTopic(threadId, postId string) string {
	util.LongPrintlnPrintln("Getting Thread viewtopic: " + threadId)

	baseDomain := o.Config.BaseUrl

	threadUrl := baseDomain + "/viewtopic?t=" + threadId + "&topic_name#" + postId
	req, err := http.NewRequest("GET", threadUrl, nil)
	util.Panic(err)

	resp, err := o.Client.Do(req)
	util.Panic(err)
	defer resp.Body.Close()
	util.PrintResponseStatus(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	util.Panic(err)

	return string(body)
}

func (o *Tool) GetPostSecrets() (string, string, error) {
	util.LongPrintlnPrintln("Getting Post Secrets: ")
	baseDomain := o.Config.BaseUrl
	postUrl := baseDomain + "/post?f=44&mode=newtopic"

	req, err := http.NewRequest("GET", postUrl, nil)
	util.Panic(err)

	resp, err := o.Client.Do(req)
	util.Panic(err)
	defer resp.Body.Close()
	util.PrintResponseStatus(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	util.Panic(err)

	secret1, secret2 := parser.GetPostSecrets(string(body))
	if secret1 == "" || secret2 == "" {
		util.LongPrintlnPrintln("ERROR: Could not get post secrets")
		err := fmt.Errorf("Could not get post secrets. Closing...")
		return "", "", err
	} else {
		util.LongPrintlnPrintln("OK: Post secrets obtained successfully")
	}

	return secret1, secret2, nil
}
