package cmd

import (
	"atcoder-gli/atcoder"
	"context"
	"fmt"
	"os"

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
	ctx := context.Background()
	ac := atcoder.NewAtCoder(ctx, sessionData.GetString("cookie"))
	name, err := ac.CheckSession()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if name != "" {
		fmt.Printf("You are logged in as %s\n", name)
	} else {
		fmt.Println("You are not logged in to AtCoder")
		os.Exit(1)
	}
}
