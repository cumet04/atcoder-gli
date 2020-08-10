package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	usage := `
Login to AtCoder with USERNAME and PASSWORD.
USERNAME and PASSWORD are optional, and they are prompted if omitted.
Some actions (ex. 'acg submit') require login beforehand, so you need to login with this command.

See also 'acg help session' for current login status.
`
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "login [USERNAME] [PASSWORD]",
			Short: "Login to AtCoder",
			Long:  strings.TrimSpace(usage),
			Args:  cobra.MaximumNArgs(2),
			Run:   cobraRun(runLogin),
		})
}

func runLogin(cmd *cobra.Command, args []string) int {
	var user string
	var pass string
	var err error
	if len(args) >= 1 {
		user = args[0]
	}
	if len(args) >= 2 {
		pass = args[1]
	}
	if user == "" {
		user, err = prompt("Username", false)
		if err != nil {
			return writeError("Prompt username failed: %s", err)
		}
	}
	if pass == "" {
		pass, err = prompt("Password", true)
		if err != nil {
			return writeError("Prompt password failed: %s", err)
		}
	}

	ac := atcoder.NewAtCoder(cmd.Context(), "")
	cookie, err := ac.Login(user, pass)
	if err != nil {
		return writeError("%s", err)
	}
	if err = saveSession(cookie); err != nil {
		panic(err)
	}

	fmt.Println("Login succeeded")
	return 0
}
