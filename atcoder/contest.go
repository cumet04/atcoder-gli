package atcoder

import "time"

// A Contest represents a contest of AtCoder; ex. abc100
type Contest struct {
	// contest stats
	ID         string        `json:"id"`    // ID (ex. "abc100")
	Title      string        `json:"title"` // Title (ex. "AtCoder Beginner Contest 100")
	URL        string        `json:"url"`   // URL (ex. "https://atcoder.jp/contests/abc100")
	StartAt    time.Time     `json:"start_at"`
	Duration   time.Duration `json:"duration"`
	Registered bool          `json:"-"` // Login user registers this contest or not

	// local info (from config)
	Script    string `json:"script"` // script file name for the tasks; basename for config.template
	SampleDir string `json:"sample_dir"`
	Language  string `json:"language"`
	Command   string `json:"command"`

	Tasks []*Task
}

// AddTask append task to its tasks
func (c *Contest) AddTask(t Task) *Task {
	t.Contest = c
	c.Tasks = append(c.Tasks, &t)
	return &t
}
