package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	usage := `
Check whether current login session is alive or not.
If session is alive, it show login user's username.

See also 'acg help login'.
`
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "session",
			Short: "Check login session status",
			Long:  strings.TrimSpace(usage),
			Run:   cobraRun(runSession),
		})
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
