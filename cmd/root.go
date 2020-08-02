package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config     Config
	configData *viper.Viper
	session    string

	rootCmd = &cobra.Command{
		Use:   "acg",
		Short: "accoder-cli on go",
		Long:  "accoder-cli on golang",
	}
)

type Config struct {
	SampleDir    string `mapstructure:"sample_dir"`
	SkeletonFile string `mapstructure:"skeleton_file"`
}

// SkeletonFilePath resolves absolute path of skeleton file.
// This regards SkeletonFile as relative path from CWD or config directory.
func (c *Config) SkeletonFilePath() string {
	if c.SkeletonFile == "" {
		return ""
	}

	if filepath.IsAbs(c.SkeletonFile) {
		return c.SkeletonFile
	}

	file1 := pathAbs(filepath.Join(configDir(), c.SkeletonFile))
	if _, err := os.Stat(file1); err == nil {
		return file1
	}

	file2 := pathAbs(filepath.Join(cwd(), c.SkeletonFile))
	if _, err := os.Stat(file2); err == nil {
		return file2
	}

	exitWithError("skeleton_file is specified but the file is not found in %s, %s", file1, file2)
	return "" // MEMO: this line is unreachable; The program exits with exitWithError
}

// Execute run rootCmd
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	c := viper.New()
	c.SetConfigName("config")
	c.SetConfigType("yml")
	c.AddConfigPath(configDir())
	c.SetDefault("sample_dir", "samples")
	c.SetDefault("skeleton_file", "")
	c.ReadInConfig()
	c.Unmarshal(&config)
	configData = c

	var err error
	session, err = readSession()
	if err != nil {
		session = ""
	}
}

func configDir() string {
	confdir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(confdir, "atcoder-gli")
}

func readSession() (string, error) {
	file := filepath.Join(configDir(), "session")
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func saveSession(cookie string) error {
	if err := os.MkdirAll(configDir(), 0755); err != nil {
		return errors.Wrapf(err, "Cannot create session directory: %s", configDir())
	}

	filename := filepath.Join(configDir(), "session")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "Cannot initialize session file: %s", filename)
	}

	_, err = file.WriteString(cookie)
	if err != nil {
		return errors.Wrapf(err, "Cannot write session to file")
	}

	return nil
}
