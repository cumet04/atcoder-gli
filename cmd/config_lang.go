package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	usage := `
Search and select preferred language.
Selected language is saved in config and used when submit code.

You can search language with keyword (prompted) and choose one from them.
Search targets are all available languages in AtCoder,
and keyword is case-insensitive.
`
	configCmd.AddCommand(
		&cobra.Command{
			Use:   "lang",
			Short: "Select preferred language for submit",
			Long:  strings.TrimSpace(usage),
			Run:   cobraRun(runLang),
		})
}

func runLang(cmd *cobra.Command, args []string) int {
	ac := atcoder.NewAtCoder(cmd.Context(), session)
	all, err := ac.ListLanguages()
	if err != nil {
		return writeError("Failed to fetch list of languages", err)
	}

	q, _ := prompt("Search", false)
	var list []atcoder.Language
	for _, l := range all {
		if strings.Contains(strings.ToUpper(l.Label), strings.ToUpper(q)) {
			list = append(list, l)
		}
	}

	prompt := promptui.Select{
		Label: "Default Language",
		Items: list,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return writeError("Failed to choose language: %s", err)
	}

	if err := config.WriteDefaultLanguage(list[index].ID); err != nil {
		return writeError("Failed to write default language to config file")
	}
	fmt.Printf("Set default language as %s (%s)\n", list[index].Label, list[index].ID)
	return 0
}
