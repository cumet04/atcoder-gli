package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gosuri/uilive"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

func init() {
	cmd := newCommand(&commandArgs{
		Use:     "submit",
		Aliases: []string{"s"},
		Run:     runSubmit,
		Short:   "Submit file to a task",
		Long: `
Submit a file as answer for a task, and wait the judge is complete.
Target file is determined by looking for a file named config's template name, in current directory.
Target task is guessed from current directory.
Language is read from config value: 'language'.
		`,
		Example: `
ex 1. run in abc100/b, template = main.rb
-> submit abc100/b/main.rb for abc100's b task
		`})
	cmd.Flags().Bool("nowait", false, "exit without waiting for judge complete")
	rootCmd.AddCommand(cmd)
}

func runSubmit(cmd *cobra.Command, args []string) int {
	task, ret := runDeterminTask()
	if ret != 0 {
		return ret
	}

	lang := config.Language()
	if lang == "" {
		return writeError("Default language is not set.\n" +
			"Retry this after set it with `config lang` command.")
	}

	path := filepath.Join(cwd(), task.Script)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return writeError("Failed to read script file: %s", err)
	}
	if len(strings.TrimSpace(string(bytes))) == 0 {
		return writeError("Script file is empty: %s", path)
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
	if submission.Judge != "AC" {
		if err := ac.FetchSubmissionDetail(submission); err != nil {
			return writeError("Failed to get judge detail: %s", err)
		}
		fmt.Println(formatJudgeCases(submission.Cases))
		fmt.Printf("See https://atcoder.jp/contests/%s/submissions/%d for detail.\n",
			task.Contest.ID, submission.ID)
	}

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
	writer.Out = colorable.NewColorableStdout()
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
		switch s.Judge {
		case "AC":
			fmt.Fprintf(writer, "Judge: %s\n", aurora.Green("AC"))
		case "WA":
			fmt.Fprintf(writer, "Judge: %s\n", aurora.Red("WA"))
		default:
			if strings.Contains(s.Judge, "/") {
				fmt.Fprintf(writer, "Judge: %s\n", s.Judge)
			} else {
				fmt.Fprintf(writer, "Judge: %s\n", aurora.Yellow(s.Judge))
			}
		}

		if interval == 0 {
			break
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	return nil
}

func formatJudgeCases(cases map[string]int) string {
	var keys []string
	for status := range cases {
		keys = append(keys, status)
	}
	sort.Strings(keys)

	var strs []string
	total := 0
	for _, key := range keys {
		count := cases[key]
		total += count
		strs = append(strs, fmt.Sprintf("%sx%d", key, count))
	}
	strs = append(strs, fmt.Sprintf("/%d", total))

	return strings.Join(strs, " ")
}
