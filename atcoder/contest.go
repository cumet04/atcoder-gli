package atcoder

// A Contest represents a contest of AtCoder; ex. abc100
type Contest struct {
	ID    string  // contest's ID (ex. "abc100")
	Title string  // contest's Title (ex. "AtCoder Beginner Contest 100")
	URL   string  // contest's URL (ex. "https://atcoder.jp/contests/abc100")
	Tasks []*Task // contest's Tasks
}

// NewContest creates a instance of Contest
func NewContest(id, title, url string) *Contest {
	return &Contest{
		ID:    id,
		Title: title,
		URL:   url,
	}
}

// AddTask append task to its tasks
func (c *Contest) AddTask(t Task) *Task {
	t.Contest = c
	c.Tasks = append(c.Tasks, &t)
	return &t
}
