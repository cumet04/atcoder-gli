package atcoder

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// A Judge represents AtCoder's judge polling API's json format
type Judge struct {
	Interval int                    `json:"Interval"` // [ms]
	Result   map[string]JudgeDetail `json:"Result"`
}

// Detail returns first judgeDetail
// In this app's usecase, a Judge must have just one judgeDetail
func (j *Judge) Detail() *JudgeDetail {
	for _, d := range j.Result {
		return &d
	}
	return nil
}

// A JudgeDetail represents a part of AtCoder's judge polling API's json format
type JudgeDetail struct {
	HTML  string `json:"Html"`
	Score string `json:"Score"`
	d     *goquery.Document
}

func (jd *JudgeDetail) tag() *goquery.Document {
	if jd.d == nil {
		// To parse correctly jd.HTML as HTML5 document, add <table> context.
		// Without this, td tag is omitted.
		h := fmt.Sprintf("<table>%s</table>", jd.HTML)
		var err error
		jd.d, err = goquery.NewDocumentFromReader(strings.NewReader(h))
		if err != nil {
			panic(err)
		}
	}
	return jd.d
}

// Result returns judge result string (ex. "AC")
func (jd *JudgeDetail) Result() string {
	return jd.tag().Find("td").First().Text()
}

// Time returns time consumption string (ex. "74 ms")
func (jd *JudgeDetail) Time() string {
	return jd.tag().Find("td").Eq(1).Text()
}

// Memory returns memory consumption string (ex. "14384 KB")
func (jd *JudgeDetail) Memory() string {
	return jd.tag().Find("td").Eq(2).Text()
}
