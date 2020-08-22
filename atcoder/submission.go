package atcoder

import "time"

// A Submission represents a submission for a problem
type Submission struct {
	ID          int       // submission's ID (ex. "15593535")
	Task        *Task     // task which the submission belongs to
	Time        int       // time consumption [ms]
	Memory      int       // memory consumption [KB]
	Judge       string    // a judge status of the submission
	SubmittedAt time.Time // submission time
}

// NewSubmission creates a instance of Submission.
// Task is assumed to be injected later.
func NewSubmission(id, time, memory int, judge string, submittedAt time.Time) *Submission {
	return &Submission{
		ID:          id,
		Task:        nil,
		Time:        time,
		Memory:      memory,
		Judge:       judge,
		SubmittedAt: submittedAt,
	}
}
