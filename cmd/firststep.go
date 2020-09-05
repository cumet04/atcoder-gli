package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"

	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

func init() {
	cmd := newCommand(&commandArgs{
		Use:   "firststep",
		Run:   runFirstStep,
		Short: "Interactive setup for first user",
		Long: `
Launch a wizard for acg's initial setup for first user.
Through the wizard, you can login to atcoder in acg and setup config with descriptions.

Same function is available 'acg login' and 'acg config wizard'.
		`})
	rootCmd.AddCommand(cmd)
}

func runFirstStep(cmd *cobra.Command, args []string) int {
	ac := atcoder.NewAtCoder(cmd.Context(), "")

	fmt.Println(aurora.Underline(aurora.Bold(
		"1. Login to atcoder. Input username/password for atcoder.jp")))
	fmt.Println("Note: Username/Password are not stored, only login session (cookie) is.")
	fmt.Println("      This step is same as 'acg login'.")
	if err := execLogin(ac, "", ""); err != nil {
		return writeError("Login sequence failed: %s", err)
	}
	fmt.Println("Login succeeded !")
	fmt.Println("")

	fmt.Println(aurora.Underline(aurora.Bold(
		"2. Setup acg config. XXX")))
	fmt.Println("Note: You can see default values / descriptions later with 'acg config doc'.")
	fmt.Println("      This step is same as 'acg config wizard'.")
	if err := execConfigWizard(ac); err != nil {
		return writeError("%s", err)
	}
	if err := config.SaveConfig(); err != nil {
		return writeError("Failed to save config to file: %s", err)
	}
	fmt.Println("Config file is created !")
	fmt.Printf("The file is located at %s\n\n", configFilePath)

	fmt.Println("Now, you are ready for use acg.")
	return 0
}
