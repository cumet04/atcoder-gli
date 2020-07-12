package cmd

import (
	"atcoder-gli/atcoder"
	"context"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "login",
	Short: "login to AtCoder",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ac := atcoder.NewAtCoder(ctx)
		ac.Login(args[0], args[1])
	},
}
