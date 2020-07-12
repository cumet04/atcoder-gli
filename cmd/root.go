package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "acg",
		Short: "accoder-cli on go",
		Long:  "accoder-cli on golang",
	}
)

// Execute run rootCmd
func Execute() error {
	return rootCmd.Execute()
}

func init() {
}
