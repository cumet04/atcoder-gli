package atcoder

// A Task represents a task in a contest; ex. abc100_a
type Task struct {
	// ContestID returns task's contest id (ex. "abc100")
	ContestID string `json:"contest_id"`
	// ID returns task's ID (ex. "abc100_a")
	ID string `json:"id"`
	// Label returns task's Label (ex. "A")
	Label string `json:"label"`
	// Title returns task's Title (ex. "Happy Birthday!")
	Title string `json:"title"`
}

// NewTask creates a instance of Task
func NewTask(contestID string, id string, label string, title string) *Task {
	return &Task{
		ContestID: contestID,
		ID:        id,
		Label:     label,
		Title:     title,
	}
}
