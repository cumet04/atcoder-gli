package atcoder

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Lang is language type used as request's Accept-Language header
type Lang int

const (
	// LangJa is language: japanese
	LangJa Lang = iota
	// LangEn is language: english
	LangEn
)

// HTTPDump is flag to dump http data (for debug)
var HTTPDump bool = false

// A Client has web client params
type Client struct {
	ctx      context.Context
	client   *http.Client
	endpoint *url.URL
	lang     Lang
	session  *http.Cookie
}

// NewClient creates new Client instance with endpoint and cookie(string)
func NewClient(ctx context.Context, endpoint string, lang Lang, cookie string) Client {
	url, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		// Stop tracing redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return Client{
		ctx,
		client,
		url,
		lang,
		parseCookie(cookie),
	}
}

func parseCookie(raw string) *http.Cookie {
	if len(raw) == 0 {
		return &http.Cookie{}
	}
	header := http.Header{}
	header.Add("Cookie", raw)
	request := http.Request{Header: header}
	return request.Cookies()[0]
}

// GetCookie returns client's current cookie string
func (c *Client) GetCookie() string {
	return c.session.String()
}

// DoGet sends GET request
func (c *Client) DoGet(spath string, expect int) (*http.Response, error) {
	return c.DoGetWithParam(spath, expect, nil)
}

// DoGetWithParam sends GET request with query param
func (c *Client) DoGetWithParam(spath string, expect int, params map[string]string) (*http.Response, error) {
	resp, err := c.doRequest("GET", spath, nil, params, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "request '[GET] %s' failed", spath)
	}
	if resp.StatusCode != expect {
		msg := fmt.Sprintf("'[GET] %s' returns unexpected status: %d", spath, resp.StatusCode)
		dumpResponse(resp)
		return nil, errors.New(msg)
	}
	return resp, nil
}

// DoFormPost sends Form POST request
func (c *Client) DoFormPost(spath string, expect int, params map[string]string) (*http.Response, error) {
	formdata := []string{}
	for k, v := range params {
		value := url.QueryEscape(v)
		formdata = append(formdata, fmt.Sprintf("%s=%s", k, value))
	}

	bodyReader := strings.NewReader(strings.Join(formdata, "&"))
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := c.doRequest("POST", spath, headers, nil, bodyReader)
	if err != nil {
		return nil, errors.Wrapf(err, "request '[POST] /%s' failed", spath)
	}
	if resp.StatusCode != expect {
		msg := fmt.Sprintf("'[POST] %s' returns unexpected status: %d", spath, resp.StatusCode)
		dumpResponse(resp)
		return nil, errors.New(msg)
	}
	return resp, nil
}

func (c *Client) doRequest(method, spath string, header map[string]string, query map[string]string, body io.Reader) (*http.Response, error) {
	url := *c.endpoint
	url.Path = path.Join(url.Path, spath)

	req, err := http.NewRequestWithContext(c.ctx, method, url.String(), body)
	if err != nil {
		panic(err)
	}

	params := req.URL.Query()
	for k, v := range query {
		params.Add(k, v)
	}
	req.URL.RawQuery = params.Encode()

	for k, v := range header {
		req.Header.Add(k, v)
	}
	if c.lang == LangJa {
		req.Header.Add("Accept-Language", "ja,en;q=0.9")
	} else {
		req.Header.Add("Accept-Language", "en,ja;q=0.9")
	}
	req.AddCookie(c.session)

	resp, err := c.client.Do(req)
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "REVEL_SESSION" {
			c.session = cookie
			break
		}
	}
	return resp, err
}

func dumpResponse(resp *http.Response) {
	if !HTTPDump {
		return
	}
	dir := filepath.Join(os.TempDir(), "atcoder-gli")
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	name := fmt.Sprintf(
		"%s%s.html",
		time.Now().Format(time.RFC3339),
		strings.ReplaceAll(resp.Request.URL.RequestURI(), "/", "-"),
	)
	path := filepath.Join(dir, name)
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("dump html to %s\n", path)
}
