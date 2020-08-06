package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "session",
			Short: "check session is active or not",
			Run:   runSession,
		})
}

func runSession(cmd *cobra.Command, args []string) {
	ac := atcoder.NewAtCoder(cmd.Context(), session)
	name, err := ac.CheckSession()
	if err != nil {
		exitWithError("%s", err)
	}

	if name != "" {
		fmt.Printf("You are logged in as %s\n", name)
	} else {
		exitWithError("You are not logged in to AtCoder")
	}
}
