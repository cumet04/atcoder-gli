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
	var url string
	if len(args) > 0 {
		url = path.Join("https://atcoder.jp/contests", args[0])
	} else {
		_, contest, err := readContestInfo("")
		if err != nil {
			exitWithError("Failed to read contest file: %s", err)
		}
		if contest == nil {
			exitWithError(
				"Cannot determin contest id.\n" +
					"Specify contest id as command arg, or run command in contest directory.",
			)
		}
		url = contest.URL
	}

	err := browser.OpenURL(url)
	if err != nil {
		exitWithError("Cannot open browser: %s", err)
	}
}
