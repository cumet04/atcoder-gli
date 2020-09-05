package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

func init() {
	cmd := newCommand(&commandArgs{
		Use:   "doc",
		Run:   runConfigDoc,
		Short: "Show config description with default values",
	})
	configCmd.AddCommand(cmd)
}

func runConfigDoc(cmd *cobra.Command, args []string) int {
	for _, param := range configDefinition() {
		name := param["name"]
		value := param["default"]
		long := param["long"]

		// MEMO: similar format to config_wizard
		head := aurora.Sprintf(aurora.Bold("* %s"), name)
		fmt.Printf("%s (default: \"%s\")\n", head, value)
		fmt.Println(" " + long)
	}

	return 0
}
