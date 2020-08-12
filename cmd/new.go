package cmd

import (
	"atcoder-gli/atcoder"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	usage := `
Create new directory for CONTEST_ID and setup directories/files.
Fetch contest info from AtCoder website and download sample test cases for tasks.

For instance, create directory tree is:
abc100/
- .contest.json
+ a/
  - main.go // if skeleton_file is set in config
  + samples/
    - sample_1.in
    - sample_1.out
    - sample_2.in
    - sample_2.out
+ b/ ...
+ c/ ...
...
`
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "new CONTEST_ID",
			Short: "Create and setup new directory for a contest",
			Long:  strings.TrimSpace(usage),
			Run:   cobraRun(runNew),
			Args:  cobra.ExactArgs(1),
		})
}

func runNew(cmd *cobra.Command, args []string) int {
	// TODO: file/dir name template
	id := args[0]

	ac := atcoder.NewAtCoder(cmd.Context(), session)
	contest, err := ac.FetchContest(id)
	if err != nil {
		return writeError("Failed to fetch contest info: %s", err)
	}
	contestDir := contest.ID
	if _, err := os.Stat(contestDir); err == nil {
		fmt.Printf("Contest directory, %s, is already exists. abort.\n", contestDir)
		return 0
	}
	if err := os.MkdirAll(contestDir, 0755); err != nil {
		return writeError("Failed to create contest directory: %s", err)
	}

	for _, t := range contest.Tasks {
		taskDir := strings.ToLower(t.Label)
		taskPath := filepath.Join(contestDir, taskDir)
		samplePath := filepath.Join(taskPath, config.SampleDir())
		if err := os.MkdirAll(samplePath, 0755); err != nil {
			return writeError("Failed to create sample directory: %s", err)
		}
		t.Directory = taskDir

		skel, err := config.SkeletonFilePath()
		if err != nil {
			return writeError("%s", err)
		}
		if skel != "" {
			filename := filepath.Base(skel)
			err := copyFile(skel, filepath.Join(taskPath, filename))
			if err != nil {
				return writeError("Failed to copy skeleton file: %s\n", err)
			}
			t.Script = filename
		}

		samples, err := ac.FetchSampleInout(t.Contest.ID, t.ID)
		if err != nil {
			return writeError("Failed to fetch task info: %s", err)
		}
		for _, s := range *samples {
			name := fmt.Sprintf("sample_%s", s.Label())
			ioutil.WriteFile(filepath.Join(samplePath, name+".in"), []byte(s.Input()), 0644)
			ioutil.WriteFile(filepath.Join(samplePath, name+".out"), []byte(s.Output()), 0644)
		}
	}

	if err := saveContestInfo(*contest, contestDir); err != nil {
		return writeError("Failed to save config: %s\n", err)
	}

	fmt.Printf("Directory for %s is ready.\n", id)
	return 0
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

// cp src dst
func copyFile(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		panic(err)
	}
	if !stat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "Cannot open file: %s", src)
	}
	defer source.Close()

	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("%s file is already exists", dst)
	}
	destination, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "Cannot create file: %s", dst)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return errors.Wrapf(err, "Cannot write file: %s", dst)
	}

	return nil
}
