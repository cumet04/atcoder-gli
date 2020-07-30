package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sessionDir  string
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

	rootCmd.PersistentFlags().StringVar(&sessionDir, "session-dir", "",
		"session file directory ($XDG_CONFIG_HOME or $HOME at default)")
}

func initConfig() {
	configData = viper.New()
	configData.SetConfigName(".atcoder-gli")
	configData.SetConfigType("json")
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configData.AddConfigPath(cwd)
	configData.AddConfigPath(filepath.Dir(cwd))

	sessionData = viper.New()
	sessionData.SetConfigName("session")
	sessionData.SetConfigType("json")
	var defaultSessionDir string
	if sessionDir != "" {
		sessionData.AddConfigPath(defaultSessionDir)
		defaultSessionDir = sessionDir
	} else {
		ucd, err := os.UserConfigDir()
		if err == nil {
			sessionData.AddConfigPath(filepath.Join(ucd, "atcoder-gli"))
			defaultSessionDir = ucd
		}
		homedir, err := os.UserHomeDir()
		if err == nil {
			sessionData.AddConfigPath(filepath.Join(homedir, "atcoder-gli"))
			if defaultSessionDir == "" {
				defaultSessionDir = homedir
			}
		}
	}
	if err := sessionData.ReadInConfig(); err != nil {
		filename := filepath.Join(defaultSessionDir, "session.json")
		_, err := os.Create(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize sesion file: %s", err)
			os.Exit(1)
		}
	}
	sessionData.SetDefault("cookie", "")
}
