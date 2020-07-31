package atcoder

// A Sample represents a pair of problem's sample input/output
type Sample struct {
	problemID string
	label     string
	input     string
	output    string
}

// NewSample creates a instance of Sample
func NewSample(problemID, label, input, output string) *Sample {
	return &Sample{
		problemID: problemID,
		label:     label,
		input:     input,
		output:    output,
	}
}

// ProblemID returns problem id that sample belongs to (ex. "abc100_a")
func (s *Sample) ProblemID() string {
	return s.problemID
}

// Label returns sample's label (ex. "1" for 入力例1/出力例1)
func (s *Sample) Label() string {
	return s.label
}

// Input returns sample's input (ex. "5 4")
func (s *Sample) Input() string {
	return s.input
}

// Output returns sample's output (ex. "Yay!")
func (s *Sample) Output() string {
	return s.output
}
