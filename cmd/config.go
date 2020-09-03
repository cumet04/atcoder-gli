package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var configCmd = newCommand(&commandArgs{
	Use:   "config",
	Run:   runConfig,
	Short: "config utility",
	Long: `
Show/Write config values from/to config file.
Run with some config options, it write the value to file.
If you run this without any options and config file, new config file is created with default values.

See 'Global Flags' for available config options.
	`})

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) int {
	if err := config.SaveConfig(); err != nil {
		return writeError("Failed to write config file :%s", err)
	}
	fmt.Print(config)
	return 0
}

// Config is interface to config data / viper
type Config struct {
	viper *viper.Viper
}

var configFilePath = filepath.Join(configDir(), "config.yml")

func configDefinition() []map[string]string {
	yml := `
- name: language
	default: ""
	short: submit code's language
	long: |
		TODO
- name: template
	default: ""
	short: template file name that is copied to task directory in 'acg new'
	long: |
		TODO
- name: command
	default: "./{{.Script}}"
	short: "command template that runs in 'acg test'"
	long: |
		TODO
- name: sample_dir
	default: tests
	short: directory name where sample in/out files are stored in
	long: |
		TODO
`
	yml = strings.ReplaceAll(yml, "\t", "  ")
	var m []map[string]string
	if err := yaml.Unmarshal([]byte(yml), &m); err != nil {
		panic(err)
	}
	return m
}

// SampleDir returns current sample_dir value of config
func (c *Config) SampleDir() string {
	return c.viper.GetString("sample_dir")
}

// Template returns current template value of config
func (c *Config) Template() string {
	return c.viper.GetString("template")
}

// Language returns current language value (ID) of config
func (c *Config) Language() string {
	raw := c.viper.GetString("language")
	return strings.Split(raw, " ")[0]
}

// Command returns command value of config
func (c *Config) Command() string {
	return c.viper.GetString("command")
}

// SaveConfig write config to default config path
func (c *Config) SaveConfig() error {
	dir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrapf(err, "Cannot create config directory: %s", dir)
	}
	return c.viper.WriteConfigAs(configFilePath)
}

// TemplateFilePath resolves absolute path of template file.
// This regards Template as relative path from CWD or config directory.
func (c *Config) TemplateFilePath() (string, error) {
	skel := c.Template()
	if skel == "" {
		return "", nil
	}

	if filepath.IsAbs(skel) {
		return skel, nil
	}

	file1 := pathAbs(filepath.Join(configDir(), skel))
	if _, err := os.Stat(file1); err == nil {
		return file1, nil
	}

	file2 := pathAbs(filepath.Join(cwd(), skel))
	if _, err := os.Stat(file2); err == nil {
		return file2, nil
	}

	return "", errors.New(fmt.Sprintf(
		"template is specified but the file is not found in %s, %s", file1, file2,
	))
}

func (c Config) String() string {
	s := c.viper.AllSettings()
	bs, err := yaml.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
