package atcoder_test

import (
	"atcoder-gli/atcoder"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckSession_LoggedIn(t *testing.T) {
	ac := atcoder.NewAtCoderTester(&StaticRoundTripper{
		pattern: "logged_in",
		baseResponse: http.Response{
			StatusCode: 200,
		},
	})
	actual, err := ac.CheckSession()
	if err != nil {
		t.Fatalf("Error on ac.CheckSession: %s", err)
	}
	if actual != "cumet04" {
		t.Fail()
	}
}

func TestCheckSession_Anonymous(t *testing.T) {
	ac := atcoder.NewAtCoderTester(&StaticRoundTripper{
		pattern: "anonymous",
		baseResponse: http.Response{
			StatusCode: 200,
		},
	})
	actual, err := ac.CheckSession()
	if err != nil {
		t.Fatalf("Error on ac.CheckSession: %s", err)
	}
	if actual != "" {
		t.Fail()
	}
}

// ----- utility codes -----

var htmlDir string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	htmlDir = filepath.Join(filepath.Dir(cwd), "test", "html")
}

// StaticRoundTripper is dummy RoundTripper that returns static file response
// corresponding to request path.
// Files are searched from test/html/ directory.
type StaticRoundTripper struct {
	pattern      string
	baseResponse http.Response
}

// RoundTrip returns static file response from request.
// The response is based on baseResponse and add Body from the file.
// This is for http.RoundTripper interface.
func (rt *StaticRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var filename string
	if rt.pattern == "" {
		filename = "index.html"
	} else {
		filename = fmt.Sprintf("index.%s.html", rt.pattern)
	}
	file, err := os.Open(filepath.Join(htmlDir, req.URL.Path, filename))
	if err != nil {
		panic(err)
	}
	resp := rt.baseResponse
	resp.Body = file

	return &resp, nil
}
