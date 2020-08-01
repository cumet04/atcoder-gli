package cmd

import (
	"atcoder-gli/atcoder"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
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
	// TODO: file/dir name template
	id := args[0]

	ac := atcoder.NewAtCoder(cmd.Context(), sessionData.GetString("cookie"))
	contest, err := ac.FetchContest(id)
	if err != nil {
		exitWithError("Failed to fetch contest info: %s", err)
	}
	contestDir := contest.ID()
	if _, err := os.Stat(contestDir); err == nil {
		fmt.Printf("Contest directory, %s, is already exists. abort.\n", contestDir)
		return
	}
	if err := os.MkdirAll(contestDir, 0755); err != nil {
		exitWithError("Failed to create contest directory: %s", err)
	}

	for _, p := range contest.Problems() {
		problemDir := filepath.Join(contestDir, strings.ToLower(p.Label()))
		sampleDir := filepath.Join(problemDir, config.SampleDir)
		if err := os.MkdirAll(sampleDir, 0755); err != nil {
			exitWithError("Failed to create sample directory: %s", err)
		}

		if config.SkeletonFile != "" {
			err := copyFile(
				filepath.Join(config.Root, config.SkeletonFile),
				filepath.Join(problemDir, filepath.Base(config.SkeletonFile)),
			)
			if err != nil {
				exitWithError("Failed to copy skeleton file: %s\n", err)
			}
		}

		samples, err := ac.FetchSampleInout(p.ContestID(), p.ID())
		if err != nil {
			exitWithError("Failed to fetch problem info: %s", err)
		}
		for _, s := range *samples {
			name := fmt.Sprintf("sample_%s", s.Label())
			ioutil.WriteFile(filepath.Join(sampleDir, name+".in"), []byte(s.Input()), 0644)
			ioutil.WriteFile(filepath.Join(sampleDir, name+".out"), []byte(s.Output()), 0644)
		}
	}

	if err := saveConfig(); err != nil {
		exitWithError("Failed to save config: %s\n", err)
	}

	fmt.Printf("Directory for %s is ready.\n", id)
}

func copyFile(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		panic(err)
	}
	if !stat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "Cannot open file: %s", src)
	}
	defer source.Close()

	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("%s file is already exists", dst)
	}
	destination, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "Cannot create file: %s", dst)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return errors.Wrapf(err, "Cannot write file: %s", dst)
	}

	return nil
}
