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

var configCmd = newCommand(&commandArgs{
	Use:   "config",
	Run:   runConfig,
	Short: "config utility",
	Long: `
Show/Write config values from/to config file.
Run with some config options, it write the value to file.
If you run this without any options and config file, new config file is created with default values.

See 'Global Flags' for available config options.
	`})

func init() {
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

var configFilePath = filepath.Join(configDir(), "config.yml")

func configDefinition() []map[string]string {
	yml := `
- name: language
	default: ""
	short: language ID specified with code submission
	long: |
		Language ID specified with code submission.
		Selectable languages are in AtCoder's submission page, and this parameter's value is a language's (AtCoder internal) ID.
		Characters after a space included in a value are ignored, so you can leave a memo such as language name.
	long_ja: |
		コード提出時に指定する言語種別のID。
		言語はAtCoderの提出ページにて選択できるもので、このパラメータの値には言語に紐付いた（AtCoder内部の）IDを指定する。
		値に含まれるスペース以降の文字は無視されるため、言語名などをメモとして残すことができる。
		'acg wizard'や'acg firststep'では'{ID} {言語名}'のフォーマットで設定される。
- name: template
	default: ""
	short: template file name that is copied to task directory in 'acg new'
	long: |
		A template file path that is copied and located at each task directory in contest directory set.
		This accepts relative path from the config file or acg command's working directory.
	long_ja: |
		'acg new'でコンテスト関連ディレクトリを作成する際、各問題ディレクトリに配置されるソースコードのテンプレートファイルのパス。
		設定ファイルもしくはコマンドの実行ディレクトリからの相対パスで記述する。
- name: command
	default: "./{{.Script}}"
	short: command template for local test in 'acg test'
	long: |
		Command template for local test in 'acg test'.
		See 'acg test --help' for detail behaviour.
		Available template value are below:
		- {{.Script}} : filename of 'template' config value without directory path
	long_ja: |
		'acg test'でローカルテストを実行するためのコマンドのテンプレート。
		動作の詳細は'acg test --help'を参照。
		テンプレート内で利用可能な値は以下:
		- {{.Script}} : 'template'設定のうち、ディレクトリを除いたファイル名の部分
- name: sample_dir
	default: tests
	short: directory name where sample in/out files are stored in
	long: |
		Directory name where sample in/out files are stored in.
		'acg new' creates directory with the name in each task directory, and sample in/out files are created in the directory.
	long_ja: |
		問題のサンプル入出力ファイルが配置されるディレクトリ名。
		'acg new'にて、問題ディレクトリの下にこの名前でディレクトリが作成され、そのディレクトリ内にサンプル入出力ファイルが作成される。
`
	// TODO: implement for 'long_ja'
	yml = strings.ReplaceAll(yml, "\t", "  ")
	var m []map[string]string
	if err := yaml.Unmarshal([]byte(yml), &m); err != nil {
		panic(err)
	}
	return m
}

// SampleDir returns current sample_dir value of config
func (c *Config) SampleDir() string {
	return c.viper.GetString("sample_dir")
}

// Template returns current template value of config
func (c *Config) Template() string {
	return c.viper.GetString("template")
}

// Language returns current language value (ID) of config
func (c *Config) Language() string {
	raw := c.viper.GetString("language")
	return strings.Split(raw, " ")[0]
}

// Command returns command value of config
func (c *Config) Command() string {
	return c.viper.GetString("command")
}

// SaveConfig write config to default config path
func (c *Config) SaveConfig() error {
	dir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrapf(err, "Cannot create config directory: %s", dir)
	}
	return c.viper.WriteConfigAs(configFilePath)
}

// TemplateFilePath resolves absolute path of template file.
// This regards Template as relative path from CWD or config directory.
func (c *Config) TemplateFilePath() (string, error) {
	skel := c.Template()
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
		"template is specified but the file is not found in %s, %s", file1, file2,
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
