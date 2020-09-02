package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd := newCommand(&commandArgs{
		Use:   "default",
		Run:   runConfigDefault,
		Short: "Show default config with description",
		Long: `
			TODO:
		`})
	configCmd.AddCommand(cmd)
}

func runConfigDefault(cmd *cobra.Command, args []string) int {
	fmt.Println("TODO: implement")

	return 0
}
