package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		&cobra.Command{
			// TODO: argsの説明（省略したらプロンプトになるなど）を入れたい
			Use:   "login [USERNAME] [PASSWORD]",
			Short: "login to AtCoder",
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
