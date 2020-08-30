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
	cmd := newCommand(&commandArgs{
		Use:     "test",
		Aliases: []string{"t"},
		Run:     runTest,
		Short:   "Run test with sample in/outs",
		Long: `
Run command with sample inputs and judge with corresponding outputs.
It uses a script and samples that are in current directory.

In default, all samples are tested in order.
		`,
		Example: `
In all cases, Current directory tree is:
+ abc100/a/
	- main.rb
	+ samples/
		- sample-1.in
		- sample-1.out
		- sample-2.in
		- sample-2.out
and config.command="ruby ./{{.ScriptFile}}", config.template="main.rb".

ex1. 'acg test'
-> run 'ruby main.rb' with stdin(sample-1.in) and judge stdout with sample-1.out.
		and same with sample-2.in / sample-2.out

ex2. 'acg test --number 1'
-> run 'ruby main.rb' with stdin(sample-1.in) and judge stdout with sample-1.out.
		it's all. sample-2 is not tested.

ex3. 'acg test --justrun --number 2'
-> run 'ruby main.rb' with stdin(sample-2.in) and show stdout/stderr of the command.
		Judge is not executed.
		`})
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

	sampleDir := filepath.Join(cwd(), task.SampleDir) // cwd should be task directory if runDeterminTask() is ok
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
		ScriptFile: task.Script,
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

		var err error
		if justrun {
			fmt.Printf("* %s\n", name)
			err = execTestJustRun(cmd.Context(), command, infile)
		} else {
			fmt.Printf("* %s ... ", name)
			err = execTestWithJudge(cmd.Context(), command, infile, outfile)
		}
		if err != nil {
			return writeError("Command execution is failed: %s", err)
		}
	}

	return 0
}

func execTestWithJudge(ctx context.Context, command, infile, outfile string) error {
	in, err := os.Open(infile)
	if err != nil {
		return errors.Wrap(err, "Failed to read in-file")
	}
	expectedBytes, err := ioutil.ReadFile(outfile)
	if err != nil {
		return errors.Wrap(err, "Failed to read out-file")
	}
	expected := string(expectedBytes)

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Stdin = in
	bytes, err := cmd.CombinedOutput()
	actual := string(bytes)
	status := cmd.ProcessState.ExitCode()
	switch status {
	case 0:
		if atcoder.Judge(actual, expected) {
			fmt.Fprintln(colorable.NewColorableStdout(), aurora.Green("AC"))
		} else {
			fmt.Fprintln(colorable.NewColorableStdout(), aurora.Red("WA"))
			fmt.Println("expected output:")
			fmt.Println(strings.TrimSuffix(expected, "\n"))
			fmt.Println("actual output:")
			fmt.Println(strings.TrimSuffix(actual, "\n"))
		}
	case -1:
		return err
	default:
		fmt.Fprint(colorable.NewColorableStdout(), aurora.Yellow("RE"))
		fmt.Printf(", with status code = %d\n", status)
		fmt.Println(actual)
	}
	return nil
}

func execTestJustRun(ctx context.Context, command, infile string) error {
	in, err := os.Open(infile)
	if err != nil {
		return errors.Wrap(err, "Failed to read in-file")
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if cmd.ProcessState.ExitCode() == -1 {
		return err
	}
	return nil
}
