package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	cmd := newCommand(&commandArgs{
		Use:   "lang",
		Run:   runLang,
		Short: "List available languages for submit",
		Long: `
List atcoder's available languages for submit.
You can also filter languages with keyword (see 'filter' flag).
		`})
	cmd.Flags().StringP("filter", "f", "", "filter keyword for list (case-insensitive)")
	cmd.Flags().Bool("no-header", false, "Don't print header")
	rootCmd.AddCommand(cmd)
}

func runLang(cmd *cobra.Command, args []string) int {
	q, err := cmd.Flags().GetString("filter")
	if err != nil {
		panic(err)
	}
	noHeader, err := cmd.Flags().GetBool("no-header")
	if err != nil {
		panic(err)
	}

	ac := atcoder.NewAtCoder(cmd.Context(), session)
	langs, err := listLanguages(ac, q)
	if err != nil {
		return writeError("%s", err)
	}

	if !noHeader {
		fmt.Println("ID\tName")
	}
	for _, lang := range langs {
		fmt.Printf("%s\t%s\n", lang.ID, lang.Label)
	}
	return 0
}

func listLanguages(ac *atcoder.AtCoder, q string) ([]atcoder.Language, error) {
	all, err := ac.ListLanguages()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch list of languages")
	}

	var list []atcoder.Language
	for _, l := range all {
		if strings.Contains(strings.ToUpper(l.Label), strings.ToUpper(q)) {
			list = append(list, l)
		}
	}

	return list, nil
}
