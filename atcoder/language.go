package atcoder

// A Language represents a pair of language id / label
type Language struct {
	ID    string // ex. "4049"
	Label string // ex. "Ruby (2.7.1)"
}

func (l Language) String() string {
	return l.Label
}
