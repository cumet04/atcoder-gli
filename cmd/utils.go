package cmd

import (
	"fmt"
	"os"
)

func exitWithError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a)
	os.Exit(1)
}
