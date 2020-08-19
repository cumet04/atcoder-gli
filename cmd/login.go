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
			Run:   runLang,
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
		return writeError("Failed to save session: %s", err)
	}

	fmt.Println("Login succeeded")
	return 0
}

func saveSession(cookie string) error {
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
