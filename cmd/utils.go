package cmd

import (
	"atcoder-gli/atcoder"
	"encoding/json"
	"io/ioutil"

	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

func exitWithError(format string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Fprintf(os.Stderr, format+"\n", a)
	} else {
		fmt.Fprintln(os.Stderr, format)
	}
	os.Exit(1)
}

func saveContestInfo(c atcoder.Contest, path string) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	filename := filepath.Join(path, ".contest.json")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	file.Write(b)
	return nil
}

func readContestInfo(basepath string) (string, *atcoder.Contest, error) {
	search := []string{
		basepath,
		filepath.Join(basepath, ".."),
		filepath.Join(basepath, "..", ".."),
	}
	for _, dir := range search {
		file := pathAbs(filepath.Join(dir, ".contest.json"))
		if _, err := os.Stat(file); err != nil {
			continue
		}

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return "", nil, err
		}
		var contest atcoder.Contest
		if err := json.Unmarshal(b, &contest); err != nil {
			return "", nil, errors.Wrap(err, "Cannot parse contest info")
		}
		for _, task := range contest.Tasks {
			task.Contest = &contest
		}
		return file, &contest, nil
	}
	return "", nil, nil
}

func prompt(label string, mask bool) (string, error) {
	var m rune
	if mask {
		m = '*'
	}
	prompt := promptui.Prompt{
		Label: label,
		Mask:  m,
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("Empty input")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	return result, nil
}

func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func pathAbs(path string) string {
	file, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return file
}
