package atcoder

// A Task represents a task in a contest; ex. abc100_a
type Task struct {
	Contest *Contest `json:"-"` // contest which the task belongs to
	ID      string   // task's ID (ex. "abc100_a")
	Label   string   // task's Label (ex. "A")
	Title   string   // task's Title (ex. "Happy Birthday!")
}

// NewTask creates a instance of Task
func NewTask(id string, label string, title string) *Task {
	return &Task{
		Contest: nil,
		ID:      id,
		Label:   label,
		Title:   title,
	}
}
