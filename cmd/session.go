package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		newCommand(&commandArgs{
			Use:   "session",
			Run:   runSession,
			Short: "Check login session status",
			Long: `
Check whether current login session is alive or not.
If session is alive, it show login user's username.

See also 'acg help login'.
			`}))
}

func runSession(cmd *cobra.Command, args []string) int {
	ac := atcoder.NewAtCoder(cmd.Context(), session)
	name, err := ac.CheckSession()
	if err != nil {
		return writeError("%s", err)
	}

	if name != "" {
		fmt.Printf("You are logged in as %s\n", name)
	} else {
		return writeError("You are not logged in to AtCoder")
	}
	return 0
}
