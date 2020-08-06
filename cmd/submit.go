package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "submit [FILE]",
			Short: "submit file to a problem",
			Run:   runSubmit,
			Args:  cobra.MaximumNArgs(1),
		})
}

func runSubmit(cmd *cobra.Command, args []string) {
	configFile, contest, err := readContestInfo("")
	if err != nil {
		exitWithError("Failed to read contest file: %s", err)
	}
	if contest == nil {
		exitWithError(
			"Cannot determin current contest.\n" +
				"Run command under contest directory.",
		)
	}
	basedir := filepath.Dir(configFile)

	var scriptDir string
	if len(args) > 0 {
		// TODO: ファイル直指定されそうなのでコメントなり分岐なり
		scriptDir = pathAbs(args[0])
	} else {
		scriptDir = cwd()
	}
	var task *atcoder.Task
	for _, t := range contest.Tasks {
		d := pathAbs(filepath.Join(basedir, t.Directory))
		if d == scriptDir {
			task = t
		}
	}
	if task == nil {
		exitWithError(
			"Cannot determin target task.\n" +
				"Run command in task's directory, or " +
				"specify target file which is located in task directory.",
		)
	}

	lang := config.DefaultLanguage
	if lang == "" {
		exitWithError("Default language is not set.\n" +
			"Retry this after set it with `lang` command.")
	}

	bytes, err := ioutil.ReadFile(filepath.Join(scriptDir, task.Script))
	if err != nil {
		exitWithError("Failed to read script file: %s", err)
	}

	ac := atcoder.NewAtCoder(cmd.Context(), session)
	submission, err := ac.Submit(task, lang, string(bytes))
	if err != nil {
		exitWithError("Failed to submit: %s", err)
	}

	fmt.Println("Code is submitted.")
	fmt.Printf("https://atcoder.jp/contests/%s/submissions/%d\n", contest.ID, submission.ID)
}
