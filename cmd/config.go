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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config utility",
	Run:   cobraRun(runConfig),
}

func init() {
	usage := `
Show/Write config values from/to config file.
Run with some config options, it write the value to file.
If you run this without any options and config file, new config file is created with default values.

See 'Global Flags' for available config options.
	`
	// definitionに長い説明書いてここのusageに出すのがよさそう（特にskeleton_fileの仕様）
	configCmd.Long = strings.TrimSpace(usage)
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

func configDefinition() []map[string]string {
	yml := `
- name: sample_dir
  default: tests
  usage: directory name where sample in/out files are stored in
- name: skeleton_file
  default: ""
  usage: skeleton file name that is copied to task directory in 'acg new'
- name: language
  default: ""
  usage: language id used as submit code's language
- name: command
  default: "./{{.ScriptFile}}"
  usage: "command template that runs in 'acg test'"
`
	var m []map[string]string
	if err := yaml.Unmarshal([]byte(yml), &m); err != nil {
		panic(err)
	}
	return m
}

// NewConfig creates a config instance with cobra params
func NewConfig(path string, cmd *cobra.Command) *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(path)
	for _, param := range configDefinition() {
		pf := rootCmd.PersistentFlags()
		pf.String(
			param["name"],
			param["default"],
			param["usage"],
		)
		v.BindPFlag(param["name"], pf.Lookup(param["name"]))
	}
	v.ReadInConfig()

	return &Config{viper: v}
}

// SampleDir returns current sample_dir value of config
func (c *Config) SampleDir() string {
	return c.viper.GetString("sample_dir")
}

// SkeletonFile returns current skeleton_file value of config
func (c *Config) SkeletonFile() string {
	return c.viper.GetString("skeleton_file")
}

// Language returns current language value of config
func (c *Config) Language() string {
	return c.viper.GetString("language")
}

// Command returns command value of config
func (c *Config) Command() string {
	return c.viper.GetString("command")
}

// WriteLanguage set id to language, and save it to file
func (c *Config) WriteLanguage(langID string) error {
	c.viper.Set("language", langID)
	return c.SaveConfig()
}

// SaveConfig write config to default config path
func (c *Config) SaveConfig() error {
	if err := os.MkdirAll(configDir(), 0755); err != nil {
		return errors.Wrapf(err, "Cannot create config directory: %s", configDir())
	}
	return c.viper.WriteConfigAs(filepath.Join(configDir(), "config.yml"))
}

// SkeletonFilePath resolves absolute path of skeleton file.
// This regards SkeletonFile as relative path from CWD or config directory.
func (c *Config) SkeletonFilePath() (string, error) {
	skel := c.SkeletonFile()
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
		"skeleton_file is specified but the file is not found in %s, %s", file1, file2,
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
