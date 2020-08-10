package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	usage := `
Show a contest summary.
Target contest is specified by CONTEST_ID, or guessed by current directory.

If you run this command in contest directory (created by 'acg new') or under it,
target contest is guessed to the directory's contest.

If CONTEST_ID is present, use it for determining target contest (current directory is not considered).
`
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "show [CONTEST_ID]",
			Short: "Show contest summary",
			Long:  strings.TrimSpace(usage),
			Run:   cobraRun(runShow),
			Args:  cobra.MaximumNArgs(1),
		})
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

	ac := atcoder.NewAtCoder(cmd.Context(), session)
	contest, err := ac.FetchContest(id)
	if err != nil {
		return writeError("%s", err)
	}

	fmt.Printf("%s (%s)\n", contest.Title, contest.ID)
	fmt.Println("-----")
	for _, p := range contest.Tasks {
		fmt.Printf("%s - %s\n", p.Label, p.Title)
	}
	return 0
}
