package cmd

import (
	"atcoder-gli/atcoder"
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	rootCmd.AddCommand(
		&cobra.Command{
			// TODO: argsの説明（省略したらプロンプトになるなど）を入れたい
			Use:   "login [username] [password]",
			Short: "login to AtCoder",
			Args:  cobra.MaximumNArgs(2),
			Run:   runLogin,
		})
}

func runLogin(cmd *cobra.Command, args []string) {
	var user string
	var pass string
	if len(args) >= 1 {
		user = args[0]
	}
	if len(args) >= 2 {
		pass = args[1]
	}
	if user == "" {
		fmt.Print("Enter Username: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		user = scanner.Text()
	}
	if pass == "" {
		fmt.Print("Enter Password: ")
		bytes, err := terminal.ReadPassword(0)
		if err != nil {
			panic(err)
		}
		fmt.Println("")
		pass = string(bytes)
	}

	ac := atcoder.NewAtCoder(cmd.Context(), "")
	cookie, err := ac.Login(user, pass)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err = saveSession(cookie); err != nil {
		panic(err)
	}

	fmt.Println("Login succeeded")
}
