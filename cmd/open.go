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
			Run:   cobraRun(runOpen),
			Args:  cobra.MaximumNArgs(1),
		})
}

func runOpen(cmd *cobra.Command, args []string) int {
	var url string
	if len(args) > 0 {
		url = path.Join("https://atcoder.jp/contests", args[0])
	} else {
		_, contest, err := readContestInfo("")
		if err != nil {
			return writeError("Failed to read contest file: %s", err)
		}
		if contest == nil {
			return writeError(
				"Cannot determin contest id.\n" +
					"Specify contest id as command arg, or run command in contest directory.",
			)
		}
		url = contest.URL
	}

	err := browser.OpenURL(url)
	if err != nil {
		return writeError("Cannot open browser: %s", err)
	}
	return 0
}
