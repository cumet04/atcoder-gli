package atcoder

// A Contest represents a contest of AtCoder; ex. abc100
type Contest struct {
	id    string
	title string
	url   string
	tasks []Task
}

// NewContest creates a instance of Contest
func NewContest(id, title, url string, tasks []Task) *Contest {
	return &Contest{
		id:    id,
		title: title,
		url:   url,
		tasks: tasks,
	}
}

// ID returns contest's id (ex. "abc100")
func (c *Contest) ID() string {
	return c.id
}

// Title returns contest's title (ex. "AtCoder Beginner Contest 100")
func (c *Contest) Title() string {
	return c.title
}

// URL returns contest's url (ex. "https://atcoder.jp/contests/abc100")
func (c *Contest) URL() string {
	return c.url
}

// Tasks returns contest's tasks
func (c *Contest) Tasks() []Task {
	return c.tasks
}
