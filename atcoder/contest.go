package atcoder

import "time"

// A Contest represents a contest of AtCoder; ex. abc100
type Contest struct {
	ID         string // ID (ex. "abc100")
	Title      string // Title (ex. "AtCoder Beginner Contest 100")
	URL        string // URL (ex. "https://atcoder.jp/contests/abc100")
	StartAt    time.Time
	Duration   time.Duration
	Tasks      []*Task
	Registered bool `json:"-"` // Login user registers this contest or not
}

// AddTask append task to its tasks
func (c *Contest) AddTask(t Task) *Task {
	t.Contest = c
	c.Tasks = append(c.Tasks, &t)
	return &t
}
