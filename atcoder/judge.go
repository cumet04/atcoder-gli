package atcoder

import "strings"

// Judge compares actual/expected output with consideration for some ambiguity
func Judge(actual, expected string) bool {
	// consider trailing newline
	if strings.HasSuffix(actual, "\n") {
		actual = strings.TrimSuffix(actual, "\n")
	}
	if strings.HasSuffix(expected, "\n") {
		expected = strings.TrimSuffix(expected, "\n")
	}

	// TODO: 小数点の精度など

	return actual == expected
}
