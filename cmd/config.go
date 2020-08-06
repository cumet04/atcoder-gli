package cmd

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	viper           *viper.Viper `mapstructure:"-"`
	SampleDir       string       `mapstructure:"sample_dir"`
	SkeletonFile    string       `mapstructure:"skeleton_file"`
	DefaultLanguage string       `mapstructure:"default_language"` // language id
}

// NewConfig creates Config instance from config file with viper
func NewConfig(path string) *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(path)

	v.SetDefault("sample_dir", "samples")
	v.SetDefault("skeleton_file", "")
	v.SetDefault("default_language", "")

	var c Config
	v.ReadInConfig()
	v.Unmarshal(&c)

	c.viper = v
	return &c
}

// WriteDefaultLanguage set id to default language, and save it to file
func (c *Config) WriteDefaultLanguage(langID string) error {
	if err := os.MkdirAll(configDir(), 0755); err != nil {
		return errors.Wrapf(err, "Cannot create config directory: %s", configDir())
	}
	c.DefaultLanguage = langID
	c.viper.Set("default_language", langID)
	return c.viper.WriteConfigAs(filepath.Join(configDir(), "config.yml"))
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
