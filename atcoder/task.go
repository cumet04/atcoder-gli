package atcoder

// A Task represents a task in a contest; ex. abc100_a
type Task struct {
	contestID string
	id        string
	label     string
	title     string
}

// NewTask creates a instance of Task
func NewTask(contestID string, id string, label string, title string) *Task {
	return &Task{
		contestID: contestID,
		id:        id,
		label:     label,
		title:     title,
	}
}

// ContestID returns task's contest id (ex. "abc100")
func (p *Task) ContestID() string {
	return p.contestID
}

// ID returns task's id (ex. "abc100_a")
func (p *Task) ID() string {
	return p.id
}

// Label returns task's label (ex. "A")
func (p *Task) Label() string {
	return p.label
}

// Title returns task's title (ex. "Happy Birthday!")
func (p *Task) Title() string {
	return p.title
}
