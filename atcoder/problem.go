package atcoder

// A Problem represents a problem in a contest; ex. abc100_a
type Problem struct {
	id    string
	label string
	name  string
}

// NewProblem creates a instance of Problem
func NewProblem(id string, label string, name string) *Problem {
	return &Problem{
		id:    id,
		label: label,
		name:  name,
	}
}

// ID returns problem's id
func (p *Problem) ID() string {
	return p.id
}

// Label returns problem's label
func (p *Problem) Label() string {
	return p.label
}

// Name returns problem's name
func (p *Problem) Name() string {
	return p.name
}
