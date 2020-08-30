package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		newCommand(&commandArgs{
			Use:   "show [CONTEST_ID]",
			Args:  cobra.MaximumNArgs(1),
			Run:   runShow,
			Short: "Show contest summary",
			Long: `
Show a contest summary.
Target contest is specified by CONTEST_ID, or guessed by current directory.

If you run this command in contest directory (created by 'acg new') or under it,
target contest is guessed to the directory's contest.

If CONTEST_ID is present, use it for determining target contest (current directory is not considered).
			`}))
}

func runShow(cmd *cobra.Command, args []string) int {
	var id string
	if len(args) > 0 {
		id = args[0]
	} else {
		_, contest, err := readContestInfo("")
		if err != nil {
			return writeError("Failed to read contest file: %s", err)
		}
		if contest == nil {
			return writeError(
				"Cannot determin contest id.\n" +
					"Specify contest id as command arg, or run command in contest directory.",
			)
		}
		id = contest.ID
	}

	// MEMO: 開催前のコンテストにshowすると/tasksが404で失敗する。
	// コマンドの挙動としてはおかしい気もするが、それで困る可能性が思い当たらないので放置
	ac := atcoder.NewAtCoder(cmd.Context(), session)
	contest, err := ac.FetchContest(id)
	if err != nil {
		return writeError("Failed to fetch contest info: %s", err)
	}

	var remaining time.Duration
	if contest.Registered {
		remaining = time.Now().Sub(contest.StartAt)
	} else {
		t, err := ac.FetchVirtualStartTime(id)
		if err != nil {
			return writeError("Failed to fetch virtual participate info: %s", err)
		}
		if t != nil && time.Now().After(*t) {
			progress := time.Now().Sub(*t)
			remaining = contest.Duration - progress
		}
	}

	fmt.Printf("%s (%s)\n", contest.Title, contest.ID)
	fmt.Println("-----")
	for _, p := range contest.Tasks {
		fmt.Printf("%s - %s\n", p.Label, p.Title)
	}
	if remaining != 0 {
		fmt.Println("-----")
		str := ""
		if remaining.Hours() > 1 {
			str = fmt.Sprintf("%02d:", int(remaining.Hours()))
		}
		if remaining.Minutes() > 1 {
			str = str + fmt.Sprintf("%02d:", int(remaining.Minutes())%60)
		}
		str = str + fmt.Sprintf("%02d", int(remaining.Seconds())%60)
		fmt.Printf("Remaining: %s\n", str)
	}
	return 0
}
