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
			Use:   "show",
			Short: "show contest summary",
			Run:   runShow,
			Args:  cobra.ExactArgs(1),
		})
}

func runShow(cmd *cobra.Command, args []string) {
	id := args[0]

	ac := atcoder.NewAtCoder(cmd.Context(), sessionData.GetString("cookie"))
	contest, err := ac.FetchContest(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("%s (%s)\n", contest.Name(), contest.ID())
	fmt.Println("-----")
	for _, p := range contest.Problems() {
		fmt.Printf("%s - %s\n", p.Label(), p.Name())
	}
}
