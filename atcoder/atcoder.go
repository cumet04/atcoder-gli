package atcoder

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// AtCoder is an application interface to get/send various data from AtCoder website
type AtCoder struct {
	client Client
}

// NewAtCoder creates new AtCoder instance with context
func NewAtCoder(ctx context.Context, cookie string) *AtCoder {
	return &AtCoder{
		NewClient(ctx, "https://atcoder.jp", LangJa, cookie),
	}
}

// Login executes login sequence with user/pass, and return cookie data
func (ac *AtCoder) Login(user, pass string) (string, error) {
	resp, err := ac.client.DoGet("/login", 200)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	token, _ := doc.Find("input[name=csrf_token]").First().Attr("value")

	resp, err = ac.client.DoFormPost("/login", 302, map[string]string{
		"username":   user,
		"password":   pass,
		"csrf_token": token,
	})
	if err != nil {
		return "", err
	}

	if resp.Header.Get("Location") != "/home" {
		msg := extractFlash(resp.Cookies(), "error")
		return "", errors.New("Login to AtCoder is failed with message: " + msg)
	}

	return ac.client.GetCookie(), nil
}

// CheckSession gets current session's user name from top page's header
func (ac *AtCoder) CheckSession() (string, error) {
	resp, err := ac.client.DoGet("/", 200)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	name := doc.Find(".header-mypage .header-mypage_btn span.bold").First().Text()

	return name, nil
}

// ListLanguages gets list of submittable languages
func (ac *AtCoder) ListLanguages() ([]Language, error) {
	resp, err := ac.client.DoGet("/contests/practice/submit", 200)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	ops := doc.Find("#select-lang-practice_1 option")

	var list []Language
	ops.Each(func(i int, s *goquery.Selection) {
		id, ok := s.Attr("value")
		if !ok {
			return
		}
		label := s.Text()
		list = append(list, Language{
			ID:    id,
			Label: html.UnescapeString(label),
		})
	})

	return list, nil
}

// FetchContest gets specified contest's summary
func (ac *AtCoder) FetchContest(id string) (*Contest, error) {
	resp, err := ac.client.DoGet(fmt.Sprintf("/contests/%s/tasks", id), 200)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	times := doc.Find(".contest-duration time")
	startAt := parseTime(times.Eq(0).Text())
	endAt := parseTime(times.Eq(1).Text())
	contest := Contest{
		ID:       id,
		Title:    doc.Find(".navbar .contest-title").First().Text(),
		URL:      resp.Request.URL.String(),
		StartAt:  startAt,
		Duration: endAt.Sub(startAt),
	}
	doc.Find("table tbody tr").Each(func(i int, tr *goquery.Selection) {
		links := tr.Find("td a")

		url, _ := links.First().Attr("href")
		dirs := strings.Split(url, "/")
		pid := dirs[len(dirs)-1]

		contest.AddTask(*NewTask(
			pid,
			links.First().Text(),
			links.Eq(1).Text(),
		))
	})

	// fetch register or not
	resp, err = ac.client.DoGet(fmt.Sprintf("/contests/%s", id), 200)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	contest.Registered = doc.Find("a[data-target='#modal-unregister']").Length() != 0

	return &contest, nil
}

// FetchSampleInout gets a task's list of sample in/out pair
func (ac *AtCoder) FetchSampleInout(contestID, taskID string) (*[]Sample, error) {
	resp, err := ac.client.DoGet(
		fmt.Sprintf("/contests/%s/tasks/%s", contestID, taskID),
		200,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// MEMO: This code doesn't work with old contests (before 2014-04)
	var samples []Sample
	pres := doc.Find("#task-statement .lang-ja h3+pre")
	for i := 0; i < pres.Length(); i += 2 {
		samples = append(samples, *NewSample(
			taskID,
			strings.Split(pres.Eq(i).Prev().Text(), " ")[1],
			strings.TrimSuffix(pres.Eq(i).Text(), "\n"),
			strings.TrimSuffix(pres.Eq(i+1).Text(), "\n"),
		))
	}
	return &samples, nil
}

// Submit posts script body as task's answer, and returns latest submission
// This method modifies task; add submission
func (ac *AtCoder) Submit(task *Task, langID, body string) (*Submission, error) {
	// post code

	path := fmt.Sprintf("/contests/%s/submit", task.Contest.ID)
	resp, err := ac.client.DoGet(path, 200)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	form := doc.Find(fmt.Sprintf("form[action='%s']", path)).First()
	token, _ := form.Find("input[name=csrf_token]").First().Attr("value")

	resp, err = ac.client.DoFormPost(path, 302, map[string]string{
		"data.TaskScreenName": task.ID,
		"data.LanguageId":     langID,
		"sourceCode":          body,
		"csrf_token":          token,
	})
	if err != nil {
		return nil, err
	}

	// get latest submission from /submissions/me

	spath := fmt.Sprintf("/contests/%s/submissions/me", task.Contest.ID)
	resp, err = ac.client.DoGet(spath, 200)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	var tds *goquery.Selection
	doc.Find("table tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		title := fmt.Sprintf("%s - %s", task.Label, task.Title)
		if tr.Find("td > a").First().Text() == title {
			tds = tr.Find("td")
			return false
		}
		return true
	})

	sidStr, _ := tds.Eq(4).Attr("data-id")
	sid, _ := strconv.Atoi(sidStr)

	at := parseTime(tds.Find("time").First().Text())
	submission := NewSubmission(
		sid,
		0, // Time/Memory are not determined at this point
		0,
		tds.Eq(6).Text(),
		at,
	)
	return task.AddSubmission(*submission), nil
}

// PollSubmissionStatus queries submission's status (judge)
// It returns zero if judge is done, or returns polling interval if not
// This method modified sub; writes judge, time, memory value
func (ac *AtCoder) PollSubmissionStatus(sub *Submission) (int, error) {
	path := fmt.Sprintf("/contests/%s/submissions/me/status/json", sub.Task.Contest.ID)
	resp, err := ac.client.DoGetWithParam(path, 200, map[string]string{
		"reload": "true",
		"sids[]": strconv.Itoa(sub.ID),
	})
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var judge JudgeStatus
	if err := json.Unmarshal(bytes, &judge); err != nil {
		return 0, err
	}

	detail := judge.Detail()
	sub.Judge = detail.Result()
	if detail.Time() != "" {
		sub.Time, _ = strconv.Atoi(strings.Split(detail.Time(), " ")[0])
	}
	if detail.Memory() != "" {
		sub.Memory, _ = strconv.Atoi(strings.Split(detail.Memory(), " ")[0])
	}

	return judge.Interval, nil
}

// FetchVirtualStartTime gets virtual participation start time of contestID's contest
// It returns nil if login user does not participate virtual contest
func (ac *AtCoder) FetchVirtualStartTime(contestID string) (*time.Time, error) {
	path := fmt.Sprintf("/contests/%s/virtual", contestID)
	resp, err := ac.client.DoGet(path, 200)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	unreg := doc.Find(fmt.Sprintf("form[action='%s/unregister']", path))
	if unreg.Length() == 0 {
		// unregister form doesn't exists -> user doesn't participate virtual contest
		return nil, nil
	}

	start := parseTime(unreg.Parent().Find("time").First().Text())
	return &start, nil
}

func extractFlash(cookies []*http.Cookie, key string) string {
	var raw string
	for _, cookie := range cookies {
		if cookie.Name == "REVEL_FLASH" {
			var err error
			raw, err = url.QueryUnescape(cookie.Value)
			if err != nil {
				panic(err)
			}
			break
		}
	}
	for _, line := range strings.Split(raw, "\x00") {
		if strings.HasPrefix(line, key+":") {
			return strings.TrimPrefix(line, key+":")
		}
	}
	return ""
}

func parseTime(str string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05-0700", str)
	if err != nil {
		panic(err)
	}
	return t
}
