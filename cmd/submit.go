package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
)

func init() {
	usage := `
Submit a file as answer for a task, and wait the judge is complete.
Target file is determined by looking for a file named config's skeleton_file name, in current directory.
Target task is guessed from current directory.
Language is read from config value: 'language'.

ex 1. run in abc100/b, skeleton_file = main.rb
-> submit abc100/b/main.rb for abc100's b task
`
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Submit file to a task",
		Run:   cobraRun(runSubmit),
		Long:  strings.TrimSpace(usage),
	}
	cmd.Flags().Bool("nowait", false, "exit without waiting for judge complete")
	rootCmd.AddCommand(cmd)
}

func runSubmit(cmd *cobra.Command, args []string) int {
	task, ret := runDeterminTask()
	if ret != 0 {
		return ret
	}
	contest := task.Contest

	lang := config.Language()
	if lang == "" {
		return writeError("Default language is not set.\n" +
			"Retry this after set it with `config lang` command.")
	}

	bytes, err := ioutil.ReadFile(filepath.Join(cwd(), task.Script))
	if err != nil {
		return writeError("Failed to read script file: %s", err)
	}

	ac := atcoder.NewAtCoder(cmd.Context(), session)
	submission, err := ac.Submit(task, lang, string(bytes))
	if err != nil {
		return writeError("Failed to submit: %s", err)
	}

	fmt.Println("Code is submitted.")
	if b, _ := cmd.Flags().GetBool("nowait"); !b {
		if err := waitForJudge(ac, submission); err != nil {
			return writeError("Error on waiting judge: %s", err)
		}
	}
	fmt.Printf("See: https://atcoder.jp/contests/%s/submissions/%d\n", contest.ID, submission.ID)

	return 0
}

func runDeterminTask() (*atcoder.Task, int) {
	configFile, contest, err := readContestInfo("")
	if err != nil {
		return nil, writeError("Failed to read contest file: %s", err)
	}
	if contest == nil {
		return nil, writeError(
			"Cannot determin current contest.\n" +
				"Run command under contest directory.",
		)
	}
	basedir := filepath.Dir(configFile)

	scriptDir := cwd()
	var task *atcoder.Task
	for _, t := range contest.Tasks {
		d := pathAbs(filepath.Join(basedir, t.Directory))
		if d == scriptDir {
			task = t
		}
	}
	if task == nil {
		return nil, writeError(
			"Cannot determin target task.\n" +
				"Run command in task's directory.",
		)
	}

	return task, 0
}

func waitForJudge(ac *atcoder.AtCoder, s *atcoder.Submission) error {
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	// initial wait; 5000ms
	for i := 0; i <= 25; i++ {
		msg := "Waiting for judge "
		for j := 0; j <= (i % 5); j++ {
			msg = msg + "."
		}
		fmt.Fprintln(writer, msg)
		time.Sleep(200 * time.Millisecond)
	}

	for {
		interval, err := ac.PollSubmissionStatus(s)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(writer, "Judge: %s\n", s.Judge)

		if interval == 0 {
			break
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	return nil
}
