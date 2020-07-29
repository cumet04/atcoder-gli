package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configDir   string
	configData  *viper.Viper
	sessionData *viper.Viper

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

	rootCmd.PersistentFlags().StringVar(&configDir, "config", "",
		"config directory (CWD at default)")
}

func initConfig() {
	var err error
	configData, err = initConfigFile("config", configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config file: %s", err)
		os.Exit(1)
	}
	sessionData, err = initConfigFile("session", configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize sesion file: %s", err)
		os.Exit(1)
	}

	sessionData.SetDefault("cookie", "")
}

func initConfigFile(name string, dir string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(name)
	v.SetConfigType("json")

	userConfigDir, err := os.UserConfigDir()
	if err == nil {
		v.AddConfigPath(filepath.Join(userConfigDir, "atcoder-gli"))
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	v.AddConfigPath(cwd)

	if dir != "" {
		v.AddConfigPath(dir)
	}

	if err := v.ReadInConfig(); err != nil {
		filename := filepath.Join(cwd, name+".json")
		_, err := os.Create(filename)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannot create file: %s", filename))
		}
	}
	return v, nil
}
