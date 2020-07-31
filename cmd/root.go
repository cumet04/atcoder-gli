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
		defaultSessionDir = sessionDir
		sessionData.AddConfigPath(defaultSessionDir)
	} else {
		ucd, err := os.UserConfigDir()
		if err == nil {
			defaultSessionDir = filepath.Join(ucd, "atcoder-gli")
			sessionData.AddConfigPath(defaultSessionDir)
		}
	}
	if err := sessionData.ReadInConfig(); err != nil {
		if err := os.MkdirAll(defaultSessionDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create default session dir: %s", err)
			os.Exit(1)
		}
		filename := filepath.Join(defaultSessionDir, "session.json")
		if _, err := os.Create(filename); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize session file: %s", err)
			os.Exit(1)
		}
	}
	sessionData.SetDefault("cookie", "")
}
