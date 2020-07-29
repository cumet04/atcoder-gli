package atcoder

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

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

// CheckSession attempt to get current session's user name from top page's header
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
	name := doc.Find(".header-mypage .user-gray").First().Text()

	return name, nil
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
