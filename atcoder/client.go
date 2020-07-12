package atcoder

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Client struct {
	endpoint  *url.URL
	client    *http.Client
	ctx       context.Context
	session   *http.Cookie
	csrfToken string
}

// NewClient creates new Client instance with endpoint and cookie(string)
func NewClient(ctx context.Context, endpoint string, cookie string) Client {
	url, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	return Client{
		url,
		http.DefaultClient,
		ctx,
		parseCookie(cookie),
		"",
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
func (c *Client) DoGet(spath string) (*http.Response, error) {
	return c.doRequest("GET", spath, nil, nil)
}

// DoFormPost sends Form POST request
func (c *Client) DoFormPost(spath string, params map[string]string) (*http.Response, error) {
	formdata := []string{}
	for k, v := range params {
		value := url.QueryEscape(v)
		formdata = append(formdata, fmt.Sprintf("%s=%s", k, value))
	}

	bodyReader := strings.NewReader(strings.Join(formdata, "&"))
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	return c.doRequest("POST", spath, headers, bodyReader)
}

func (c *Client) doRequest(method, spath string, header map[string]string, body io.Reader) (*http.Response, error) {
	url := *c.endpoint
	url.Path = path.Join(url.Path, spath)

	req, err := http.NewRequestWithContext(c.ctx, method, url.String(), body)
	if err != nil {
		panic(err)
	}
	for k, v := range header {
		req.Header.Add(k, v)
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
