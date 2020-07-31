package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "new",
			Short: "create files for a contest",
			Long:  "create new directories & files for a specified contest",
			Run:   runNew,
			Args:  cobra.ExactArgs(1),
		})
}

func runNew(cmd *cobra.Command, args []string) {
	// TODO: put config file, file name template
	// put skeleton program file
	id := args[0]
	sampleDir := "tests"

	ac := atcoder.NewAtCoder(cmd.Context(), sessionData.GetString("cookie"))
	contest, err := ac.FetchContest(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	contestDir := contest.ID()
	if err := os.MkdirAll(contestDir, 0755); err != nil {
		panic(err)
	}

	for _, p := range contest.Problems() {
		sampleDir := filepath.Join(contestDir, strings.ToLower(p.Label()), sampleDir)
		if err := os.MkdirAll(sampleDir, 0755); err != nil {
			panic(err)
		}

		samples, err := ac.FetchSampleInout(p.ContestID(), p.ID())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		for _, s := range *samples {
			name := fmt.Sprintf("sample_%s", s.Label())
			ioutil.WriteFile(filepath.Join(sampleDir, name+".in"), []byte(s.Input()), 0644)
			ioutil.WriteFile(filepath.Join(sampleDir, name+".out"), []byte(s.Output()), 0644)
		}
	}
}
