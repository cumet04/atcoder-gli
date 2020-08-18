package cmd

import (
	"atcoder-gli/atcoder"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	usage := `
Run command with sample inputs and judge with corresponding outputs.
It uses a script and samples that are in current directory.

In default, all samples are tested in order.
`
	example := `  In all cases, Current directory tree is:
  + abc100/a/
    - main.rb
    + samples/
      - sample-1.in
      - sample-1.out
      - sample-2.in
      - sample-2.out
  and config.command="ruby ./{{.ScriptFile}}", config.skeleton_file="main.rb".
  
  ex1. 'acg test'
  -> run 'ruby main.rb' with stdin(sample-1.in) and judge stdout with sample-1.out.
     and same with sample-2.in / sample-2.out
  
  ex2. 'acg test --number 1'
  -> run 'ruby main.rb' with stdin(sample-1.in) and judge stdout with sample-1.out.
     it's all. sample-2 is not tested.
  
  ex3. 'acg test --justrun --number 2'
  -> run 'ruby main.rb' with stdin(sample-2.in) and show stdout/stderr of the command.
     Judge is not executed.`

	cmd := &cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
		Short:   "Run test with sample in/outs",
		Long:    strings.TrimSpace(usage),
		Example: example,
		Run:     cobraRun(runTest),
	}
	cmd.Flags().BoolP("justrun", "r", false, "just run, without judge")
	cmd.Flags().StringP("number", "n", "", "test only specified number; set '1' for 'sample-1.in/out'")
	rootCmd.AddCommand(cmd)
}

type commandEnv struct {
	ScriptFile string
}

type testResult struct {
	Judge  string
	Output string
	Status int
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
		ScriptFile: filepath.Base(config.SkeletonFile()),
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

		fmt.Printf("* %s ... ", name)
		result, err := execTestRun(cmd.Context(), command, infile, outfile)
		if err != nil {
			return writeError("Command execution is failed: %s", err)
		}

		if justrun {
			fmt.Println("Done")
			fmt.Print(result.Output)
			continue
		}

		switch result.Judge {
		case "AC":
			fmt.Fprintln(colorable.NewColorableStdout(), aurora.Green("AC"))
		case "WA":
			fmt.Fprintln(colorable.NewColorableStdout(), aurora.Red("WA"))
			fmt.Println("expected output:")
			expected, err := ioutil.ReadFile(outfile)
			if err != nil {
				return writeError("Failed to read out-file: %s", err)
			}
			fmt.Println(string(expected))
		case "RE":
			fmt.Fprint(colorable.NewColorableStdout(), aurora.Yellow("RE"))
			fmt.Printf(", with status code = %d\n", result.Status)
			fmt.Println(result.Output)
		}
	}

	return 0
}

func execTestRun(ctx context.Context, command, infile, outfile string) (*testResult, error) {
	in, err := os.Open(infile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read in-file")
	}
	expected, err := ioutil.ReadFile(outfile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read out-file")
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Stdin = in
	bytes, err := cmd.CombinedOutput()
	actual := string(bytes)
	status := cmd.ProcessState.ExitCode()
	if err != nil && status != -1 {
		return &testResult{
			Judge:  "RE",
			Output: actual,
			Status: status,
		}, nil
	}

	var judge string
	if atcoder.Judge(actual, string(expected)) {
		judge = "AC"
	} else {
		judge = "WA"
	}
	return &testResult{
		Judge:  judge,
		Output: actual,
		Status: status,
	}, nil
}
