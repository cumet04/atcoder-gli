package cmd

import (
	"path"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "open [CONTEST_ID]",
			Short: "open contest page with browser",
			Run:   runOpen,
			Args:  cobra.MaximumNArgs(1),
		})
}

func runOpen(cmd *cobra.Command, args []string) {
	var id string
	if len(args) > 0 {
		id = args[0]
	} else {
		id := currentContestDir()
		if id == "" {
			exitWithError(
				"Cannot determin contest id.\n" +
					"Specify contest id as command arg, or run command in contest directory.",
			)
		}
		// TODO: read url from config
	}
	err := browser.OpenURL(path.Join("https://atcoder.jp/contests", id))
	if err != nil {
		exitWithError("Cannot open browser: %s", err)
	}
}
