package atcoder

// A Problem represents a problem in a contest; ex. abc100_a
type Problem struct {
	contestID string
	id        string
	label     string
	name      string
}

// NewProblem creates a instance of Problem
func NewProblem(contestID string, id string, label string, name string) *Problem {
	return &Problem{
		contestID: contestID,
		id:        id,
		label:     label,
		name:      name,
	}
}

// ContestID returns problem's contest id (ex. "abc100")
func (p *Problem) ContestID() string {
	return p.contestID
}

// ID returns problem's id (ex. "abc100_a")
func (p *Problem) ID() string {
	return p.id
}

// Label returns problem's label (ex. "A")
func (p *Problem) Label() string {
	return p.label
}

// Name returns problem's name (ex. "Happy Birthday!")
func (p *Problem) Name() string {
	return p.name
}
