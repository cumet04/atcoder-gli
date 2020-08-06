package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	config  Config
	session string

	rootCmd = &cobra.Command{
		Use:   "acg",
		Short: "accoder-cli on go",
		Long:  "accoder-cli on golang",
	}
)

// Execute run rootCmd
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	config = *NewConfig(configDir())

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
