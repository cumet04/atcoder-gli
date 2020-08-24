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
		t.Errorf("\nexpected: %v\nactual  : %v", "cumet04", actual)
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
		t.Errorf("\nexpected: %v\nactual  : %v", "", actual)
	}
}

func TestFetchSubmissionDetail_OnContest(t *testing.T) {
	ac := atcoder.NewAtCoderTester(&StaticRoundTripper{
		pattern: "on_contest",
		baseResponse: http.Response{
			StatusCode: 200,
		},
	})

	sub := (&atcoder.Contest{ID: "abc176"}).
		AddTask(
			atcoder.Task{ID: "abc176_d"},
		).AddSubmission(
		atcoder.Submission{ID: 16148321},
	)
	err := ac.FetchSubmissionDetail(sub)
	if err != nil {
		t.Fatalf("Error on ac.FetchSubmissionDetail: %s", err)
	}

	// MEMO: Maps are now printed in key-sorted order to ease testing.
	// https://stackoverflow.com/questions/18208394/how-to-test-the-equivalence-of-maps-in-golang/54173309#54173309
	expected := map[string]int{
		"AC":  8,
		"WA":  5,
		"TLE": 1,
	}
	if fmt.Sprint(sub.Cases) != fmt.Sprint(expected) {
		t.Errorf("\nexpected: %v\nactual  : %v", expected, sub.Cases)
	}
}

func TestFetchSubmissionDetail_After(t *testing.T) {
	ac := atcoder.NewAtCoderTester(&StaticRoundTripper{
		pattern: "after",
		baseResponse: http.Response{
			StatusCode: 200,
		},
	})

	sub := (&atcoder.Contest{ID: "abc176"}).
		AddTask(
			atcoder.Task{ID: "abc176_d"},
		).AddSubmission(
		atcoder.Submission{ID: 16148321},
	)
	err := ac.FetchSubmissionDetail(sub)
	if err != nil {
		t.Fatalf("Error on ac.FetchSubmissionDetail: %s", err)
	}

	expected := map[string]int{
		"AC":  8,
		"WA":  5,
		"TLE": 1,
	}
	if fmt.Sprint(sub.Cases) != fmt.Sprint(expected) {
		t.Errorf("\nexpected: %v\nactual  : %v", expected, sub.Cases)
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
