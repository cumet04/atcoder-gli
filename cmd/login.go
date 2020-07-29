package cmd

import (
	"atcoder-gli/atcoder"
	"context"
	"fmt"
	"os"

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
		cookie, err := ac.Login(args[0], args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		sessionData.Set("cookie", cookie)
		err = sessionData.WriteConfig()
		if err != nil {
			panic(err)
		}

		fmt.Println("Login succeeded")
	},
}
