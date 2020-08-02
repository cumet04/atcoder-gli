package atcoder

// A Contest represents a contest of AtCoder; ex. abc100
type Contest struct {
	// ID returns contest's ID (ex. "abc100")
	ID string `json:"id"`
	// Title returns contest's Title (ex. "AtCoder Beginner Contest 100")
	Title string `json:"title"`
	// URL returns contest's URL (ex. "https://atcoder.jp/contests/abc100")
	URL string `json:"url"`
	// Tasks returns contest's Tasks
	Tasks []Task `json:"tasks"`
}

// NewContest creates a instance of Contest
func NewContest(id, title, url string, tasks []Task) *Contest {
	return &Contest{
		ID:    id,
		Title: title,
		URL:   url,
		Tasks: tasks,
	}
}
