package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		newCommand(&commandArgs{
			Use:   "login [USERNAME] [PASSWORD]",
			Args:  cobra.MaximumNArgs(2),
			Run:   runLogin,
			Short: "Login to AtCoder",
			Long: `
Login to AtCoder with USERNAME and PASSWORD.
USERNAME and PASSWORD are optional, and they are prompted if omitted.
Some actions (ex. 'acg submit') require login beforehand, so you need to login with this command.

See also 'acg help session' for current login status.
			`}))
}

func runLogin(cmd *cobra.Command, args []string) int {
	var user string
	var pass string
	if len(args) >= 1 {
		user = args[0]
	}
	if len(args) >= 2 {
		pass = args[1]
	}

	ac := atcoder.NewAtCoder(cmd.Context(), "")
	if err := execLogin(ac, user, pass); err != nil {
		return writeError("Login sequence failed: %s", err)
	}

	fmt.Println("Login succeeded")
	return 0
}

func execLogin(ac *atcoder.AtCoder, user, pass string) error {
	if user == "" {
		user = prompt(promptParam{Label: "Username"})
	}
	if pass == "" {
		pass = prompt(promptParam{Label: "Password", Mask: '*'})
	}

	cookie, err := ac.Login(user, pass)
	if err != nil {
		return err
	}

	// save session
	if err := os.MkdirAll(configDir(), 0755); err != nil {
		return errors.Wrapf(err, "Cannot create session directory: %s", configDir())
	}

	filename := filepath.Join(configDir(), "session")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "Cannot initialize session file: %s", filename)
	}

	_, err = file.WriteString(cookie)
	if err != nil {
		return errors.Wrapf(err, "Cannot write session to file")
	}

	return nil
}
