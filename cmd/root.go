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
	root         string
	SampleDir    string `mapstructure:"sample_dir"`
	SkeletonFile string `mapstructure:"skeleton_file"`
}

// Root returns resolved absolute path of config.root
func (c *Config) Root() string {
	if filepath.IsAbs(c.root) {
		return c.root
	}
	confdir := filepath.Dir(configData.ConfigFileUsed())
	dir, err := filepath.Abs(filepath.Join(confdir, c.root))
	if err != nil {
		panic(err)
	}
	return dir
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
	c.SetConfigName(".atcoder-gli")
	c.SetConfigType("json")
	c.AddConfigPath(cwd())
	c.AddConfigPath(filepath.Dir(cwd()))
	c.AddConfigPath(filepath.Dir(filepath.Dir(cwd())))
	c.SetDefault("root", ".")
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

func saveConfig() error {
	file := filepath.Join(config.Root(), ".atcoder-gli.json")
	return configData.WriteConfigAs(file)
}

// session ----------

func sessionFile() string {
	confdir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(confdir, "atcoder-gli", "session")
}

func readSession() (string, error) {
	b, err := ioutil.ReadFile(sessionFile())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func saveSession(cookie string) error {
	dir := filepath.Dir(sessionFile())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrapf(err, "Cannot create session directory: %s", dir)
	}

	file, err := os.OpenFile(sessionFile(), os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrapf(err, "Cannot initialize session file: %s", sessionFile())
	}

	_, err = file.WriteString(cookie)
	if err != nil {
		return errors.Wrapf(err, "Cannot write session to file")
	}

	return nil
}
