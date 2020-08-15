package cmd

import (
	"atcoder-gli/atcoder"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	usage := `
Submit a file as answer for a task, and wait the judge is complete.
Target file is determined by looking for a file named config's skeleton_file name, in current directory.
Target task is guessed from current directory.
Language is read from config value: 'language'.

ex 1. run in abc100/b, skeleton_file = main.rb
-> submit abc100/b/main.rb for abc100's b task
`
	cmd := &cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
		Short:   "Run test with sample in/outs",
		Run:     cobraRun(runTest),
		Long:    strings.TrimSpace(usage),
	}
	cmd.Flags().BoolP("justrun", "r", false, "just run, without judge")
	cmd.Flags().StringP("number", "n", "", "test only specified number; set '1' for 'sample-1.in/out'")
	rootCmd.AddCommand(cmd)
}

type commandEnv struct {
	ScriptFile string
}

func runTest(cmd *cobra.Command, args []string) int {
	justrun, err := cmd.Flags().GetBool("justrun")
	if err != nil {
		panic(err)
	}
	num, err := cmd.Flags().GetString("number")
	if err != nil {
		panic(err)
	}

	task, ret := runDeterminTask()
	if ret != 0 {
		return ret
	}

	// cwd should be task directory if runDeterminTask() is ok
	sampleDir := filepath.Join(cwd(), task.SampleDir)
	files, err := ioutil.ReadDir(sampleDir)
	if err != nil {
		return writeError("Cannot read sample dir: %s", err)
	}
	names := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		n := file.Name()
		if strings.HasSuffix(n, ".in") {
			base := strings.TrimSuffix(n, ".in")
			names = append(names, base)
		}
	}

	// generate command string from template
	cenv := commandEnv{
		ScriptFile: config.SkeletonFile(),
	}
	tmpl, err := template.New("command").Parse(config.Command())
	if err != nil {
		return writeError("Command template is not parsable: %s", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cenv); err != nil {
		return writeError("Cannot construct command from template: %s", err)
	}
	command := buf.String()

	// run test for samples
	for _, name := range names {
		if num != "" && fmt.Sprintf("sample-%s", num) != name {
			continue
		}

		full := filepath.Join(sampleDir, name)
		infile := full + ".in"
		outfile := full + ".out"
		if _, err := os.Stat(outfile); err != nil {
			return writeError("Output file corresponding %s is not found: %s", name+".in", outfile)
		}

		// TODO: 適当に色つけたい
		fmt.Printf("### %s\n", name)
		ok, status, err := execTestRun(cmd.Context(), command, infile, outfile)
		if err != nil {
			return writeError("Command execution is failed: %s", err)
		}
		fmt.Println("") // write \n for actual output without trailing \n

		if justrun {
			continue
		}

		if status != 0 {
			fmt.Printf("=> RE; status code = %d\n", status)
		} else if ok {
			fmt.Println("=> AC")
		} else {
			fmt.Println("=> WA")
			fmt.Println("expected output:")
			expected, err := ioutil.ReadFile(outfile)
			if err != nil {
				return writeError("Failed to read out-file: %s", err)
			}
			fmt.Println(string(expected))
		}
		fmt.Println("")
	}

	return 0
}

func execTestRun(ctx context.Context, command, infile, outfile string) (bool, int, error) {
	in, err := os.Open(infile)
	if err != nil {
		return false, -1, errors.Wrap(err, "Failed to read in-file")
	}
	expected, err := ioutil.ReadFile(outfile)
	if err != nil {
		return false, -1, errors.Wrap(err, "Failed to read out-file")
	}

	var buf bytes.Buffer
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Stdin = in
	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	status := cmd.ProcessState.ExitCode()
	if err != nil && status != -1 {
		return false, 1, errors.Wrap(err, "Failed to start command")
	}

	actual := string(buf.Bytes())
	return atcoder.Judge(actual, string(expected)), cmd.ProcessState.ExitCode(), nil
}
