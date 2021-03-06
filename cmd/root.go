package cmd

import (
	"atcoder-gli/atcoder"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	session        string
	packageVersion string
	packageCommit  string
	config         = Config{
		viper: viper.New(),
	}

	rootCmd = &cobra.Command{
		Use:     "acg",
		Short:   "accoder-cli on go",
		Version: fmt.Sprintf("%s (%s)", packageVersion, packageCommit),
	}
)

type commandArgs struct {
	Use     string
	Aliases []string
	Short   string
	Long    string
	Example string
	Run     func(cmd *cobra.Command, args []string) int
	Args    cobra.PositionalArgs
}

func newCommand(arg *commandArgs) *cobra.Command {
	indent := func(raw string) string {
		if raw == "" {
			return ""
		}
		out := ""
		for _, l := range strings.Split(strings.TrimSpace(raw), "\n") {
			out = out + "  " + l + "\n"
		}
		return strings.ReplaceAll(strings.TrimSuffix(out, "\n"), "\t", "  ")
	}
	return &cobra.Command{
		Use:     arg.Use,
		Aliases: arg.Aliases,
		Short:   arg.Short,
		Long:    strings.TrimSpace(arg.Long),
		Example: indent(arg.Example),
		Run: func(cmd *cobra.Command, args []string) {
			r := arg.Run(cmd, args)
			os.Exit(r)
		},
		Args: arg.Args,
	}
}

// Execute run rootCmd
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	config.viper.SetConfigName("config")
	config.viper.SetConfigType("yml")
	config.viper.AddConfigPath(configDir())
	config.viper.ReadInConfig()

	_, ok := os.LookupEnv("ATCODER_GLI_HTTP_DUMP")
	atcoder.HTTPDump = ok
}

func initConfig() {
	var err error
	file := filepath.Join(configDir(), "session")
	b, err := ioutil.ReadFile(file)
	if err == nil {
		session = string(b)
	} else {
		session = ""
	}
}

func writeError(format string, a ...interface{}) int {
	if len(a) > 0 {
		fmt.Fprintf(os.Stderr, format+"\n", a...)
	} else {
		fmt.Fprintln(os.Stderr, format)
	}
	return 1
}

func readContestInfo(basepath string) (string, *atcoder.Contest, error) {
	// Seach .contest.json in cwd, cwd/.., cwd/../..
	// These assume that command runs in contest root, task dir, sample dir
	search := []string{
		basepath,
		filepath.Join(basepath, ".."),
		filepath.Join(basepath, "..", ".."),
	}
	for _, dir := range search {
		file := pathAbs(filepath.Join(dir, ".contest.json"))
		if _, err := os.Stat(file); err != nil {
			continue
		}

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return "", nil, err
		}
		var contest atcoder.Contest
		if err := json.Unmarshal(b, &contest); err != nil {
			return "", nil, errors.Wrap(err, "Cannot parse contest info")
		}
		for _, task := range contest.Tasks {
			task.Contest = &contest
		}
		return file, &contest, nil
	}
	return "", nil, nil
}

type promptParam struct {
	Label   string
	Default string
	Mask    rune
}

func prompt(p promptParam) string {
	prompt := promptui.Prompt{
		Label:   p.Label,
		Mask:    p.Mask,
		Default: p.Default,
	}

	result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	return result
}

func configDir() string {
	confdir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(confdir, "atcoder-gli")
}

func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func pathAbs(path string) string {
	file, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return file
}
