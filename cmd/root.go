package cmd

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config      Config
	configData  *viper.Viper
	sessionData *viper.Viper

	rootCmd = &cobra.Command{
		Use:   "acg",
		Short: "accoder-cli on go",
		Long:  "accoder-cli on golang",
	}
)

type Config struct {
	Root         string
	SampleDir    string `mapstructure:"sample_dir"`
	SkeletonFile string `mapstructure:"skeleton_file"`
}

// Execute run rootCmd
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	c := viper.New()
	c.SetConfigName(".atcoder-gli")
	c.SetConfigType("json")
	c.AddConfigPath(cwd)
	c.AddConfigPath(filepath.Dir(cwd))
	c.SetDefault("root", ".")
	c.SetDefault("sample_dir", "samples")
	c.SetDefault("skeleton_file", "")
	c.ReadInConfig()
	c.Unmarshal(&config)
	configData = c

	s := viper.New()
	s.SetConfigName("session")
	s.SetConfigType("json")
	s.AddConfigPath(sessionDir())
	s.SetDefault("cookie", "")
	s.ReadInConfig()
	sessionData = s
}

func saveConfig() error {
	file := filepath.Join(config.Root, ".atcoder-gli.json")
	return configData.WriteConfigAs(file)
}

func sessionDir() string {
	confdir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(confdir, "atcoder-gli")
}

func saveSession(cookie string) error {
	dir := sessionDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrapf(err, "Cannot create session directory: %s", dir)
	}
	filename := filepath.Join(dir, "session.json")
	if _, err := os.OpenFile(filename, os.O_CREATE, 0644); err != nil {
		return errors.Wrapf(err, "Cannot initialize session file: %s", filename)
	}

	sessionData.Set("cookie", cookie)
	sessionData.WriteConfig()
	return nil
}
