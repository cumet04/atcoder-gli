package atcoder

// A Contest represents a contest of AtCoder; ex. abc100
type Contest struct {
	id       string
	name     string
	problems []Problem
}

// NewContest creates a instance of Contest
func NewContest(id string, name string, problems []Problem) *Contest {
	return &Contest{
		id:       id,
		name:     name,
		problems: problems,
	}
}

// ID returns contest's id (ex. "abc100")
func (c *Contest) ID() string {
	return c.id
}

// Name returns contest's name (ex. "AtCoder Beginner Contest 100")
func (c *Contest) Name() string {
	return c.name
}

// Problems returns contest's problems
func (c *Contest) Problems() []Problem {
	return c.problems
}
