package cmd

import (
	"fmt"
	"os"
)

func exitWithError(format string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Fprintf(os.Stderr, format+"\n", a)
	} else {
		fmt.Fprintln(os.Stderr, format)
	}
	os.Exit(1)
}
