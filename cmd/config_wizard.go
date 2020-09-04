package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"

	"github.com/logrusorgru/aurora/v3"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	cmd := newCommand(&commandArgs{
		Use:   "wizard",
		Run:   runWizard,
		Short: "Making config wizard",
		Long: `
			Launch wizard for making config file, and interactively setup config parameters.
		`})
	configCmd.AddCommand(cmd)
}

func runWizard(cmd *cobra.Command, args []string) int {
	ac := atcoder.NewAtCoder(cmd.Context(), session)
	err := execConfigWizard(ac)
	if err != nil {
		return writeError("%s", err)
	}
	if err := config.SaveConfig(); err != nil {
		return writeError("Failed to save config to file: %s", err)
	}
	fmt.Printf("Config file is saved at %s\n", configFilePath)

	return 0
}

func execConfigWizard(ac *atcoder.AtCoder) error {
	for _, param := range configDefinition() {
		name := param["name"]
		current := config.viper.GetString(name)
		long := param["long"]
		fmt.Println(aurora.Sprintf(aurora.Bold("* %s"), name))
		fmt.Println(" " + long)

		var value string
		switch name {
		case "language":
			q := prompt(promptParam{Label: "language search keyword (empty to skip)"})
			if q == "" {
				value = ""
			} else {
				list, err := listLanguages(ac, q)
				if err != nil {
					return err
				}

				prompt := promptui.Select{
					Label: "Select language",
					Items: list,
				}
				index, _, err := prompt.Run()
				if err != nil {
					panic(err)
				}
				value = fmt.Sprintf("%s %s", list[index].ID, list[index].Label)
			}
		default:
			value = prompt(promptParam{
				Label:   fmt.Sprintf("Input (%s)", name),
				Default: current,
			})
		}
		fmt.Println()

		config.viper.Set(name, value)
	}
	return nil
}
