package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func exitWithError(format string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Fprintf(os.Stderr, format+"\n", a)
	} else {
		fmt.Fprintln(os.Stderr, format)
	}
	os.Exit(1)
}

func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

// currentContestDir returns guessed current directory's contest id from CWD
func currentContestDir() string {
	if config.Root() != filepath.Dir(cwd()) {
		return ""
	}
	dir, err := filepath.Rel(config.Root(), cwd())
	if err != nil {
		panic(err)
	}
	return filepath.Base(dir)
}
