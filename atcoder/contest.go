package atcoder

// A Contest represents a contest of AtCoder; ex. abc100
type Contest struct {
	id    string
	title string
	tasks []Task
}

// NewContest creates a instance of Contest
func NewContest(id string, title string, tasks []Task) *Contest {
	return &Contest{
		id:    id,
		title: title,
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

// Tasks returns contest's tasks
func (c *Contest) Tasks() []Task {
	return c.tasks
}
