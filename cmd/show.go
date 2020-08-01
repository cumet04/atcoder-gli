package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "show [CONTEST_ID]",
			Short: "show contest summary",
			Run:   runShow,
			Args:  cobra.MaximumNArgs(1),
		})
}

func runShow(cmd *cobra.Command, args []string) {
	var id string
	if len(args) > 0 {
		id = args[0]
	} else {
		id = currentContestDir()
		if id == "" {
			exitWithError(
				"Cannot determin contest id.\n" +
					"Specify contest id as command arg, or run command in contest directory.",
			)
		}
	}

	ac := atcoder.NewAtCoder(cmd.Context(), session)
	contest, err := ac.FetchContest(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("%s (%s)\n", contest.Title(), contest.ID())
	fmt.Println("-----")
	for _, p := range contest.Tasks() {
		fmt.Printf("%s - %s\n", p.Label(), p.Title())
	}
}
