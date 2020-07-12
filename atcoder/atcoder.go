package atcoder

import (
	"context"

	"github.com/PuerkitoBio/goquery"
)

type AtCoder struct {
	client Client
}

func NewAtCoder(ctx context.Context) *AtCoder {
	return &AtCoder{
		NewClient(ctx, "https://atcoder.jp", ""),
	}
}

func (ac *AtCoder) Login(user, pass string) {
	resp, err := ac.client.DoGet("/login")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	token, _ := doc.Find("input[name=csrf_token]").First().Attr("value")

	resp, err = ac.client.DoFormPost("/login", map[string]string{
		"username":   user,
		"password":   pass,
		"csrf_token": token,
	})
	if err != nil {
		panic(err)
	}

	// TODO: save session
}
