package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	viper           *viper.Viper `mapstructure:"-"`
	SampleDir       string       `mapstructure:"sample_dir"`
	SkeletonFile    string       `mapstructure:"skeleton_file"`
	DefaultLanguage string       `mapstructure:"default_language"` // language id
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "[WIP] config utility",
	Run:   cobraRun(runConfig),
}

func init() {
	// 	usage := `
	// Show/Write config values from/to config file.
	// Without any options it shows current values, or with config options it write the value to file.

	// See 'acg help' for available config options.
	// `
	// 	configCmd.Long = strings.TrimSpace(usage)
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) int {
	fmt.Println(config)
	return 0
}

// NewConfig creates Config instance with viper
// if read is true, load values from config file (else values are default)
func NewConfig(path string, read bool) *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(path)

	v.SetDefault("sample_dir", "samples")
	v.SetDefault("skeleton_file", "")
	v.SetDefault("default_language", "")

	var c Config
	if read {
		v.ReadInConfig()
	}
	v.Unmarshal(&c)

	c.viper = v
	return &c
}

func (c Config) String() string {
	s := c.viper.AllSettings()
	bs, err := yaml.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

// SaveConfig write config to default config path
func (c *Config) SaveConfig() error {
	if err := os.MkdirAll(configDir(), 0755); err != nil {
		return errors.Wrapf(err, "Cannot create config directory: %s", configDir())
	}
	return c.viper.WriteConfigAs(filepath.Join(configDir(), "config.yml"))
}

// WriteDefaultLanguage set id to default language, and save it to file
func (c *Config) WriteDefaultLanguage(langID string) error {
	c.DefaultLanguage = langID
	c.viper.Set("default_language", langID)
	return c.SaveConfig()
}

// SkeletonFilePath resolves absolute path of skeleton file.
// This regards SkeletonFile as relative path from CWD or config directory.
func (c *Config) SkeletonFilePath() (string, error) {
	if c.SkeletonFile == "" {
		return "", nil
	}

	if filepath.IsAbs(c.SkeletonFile) {
		return c.SkeletonFile, nil
	}

	file1 := pathAbs(filepath.Join(configDir(), c.SkeletonFile))
	if _, err := os.Stat(file1); err == nil {
		return file1, nil
	}

	file2 := pathAbs(filepath.Join(cwd(), c.SkeletonFile))
	if _, err := os.Stat(file2); err == nil {
		return file2, nil
	}

	return "", errors.New(fmt.Sprintf(
		"skeleton_file is specified but the file is not found in %s, %s", file1, file2,
	))
}
