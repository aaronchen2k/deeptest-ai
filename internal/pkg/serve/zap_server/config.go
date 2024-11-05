package zap_server

import (
	"encoding/json"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	"strconv"

	"github.com/spf13/viper"
)

var CONFIG = Zap{
	Level:         "debug",
	Format:        "console",
	Prefix:        "[DEEPTEST]",
	Director:      "logs",
	LinkName:      "latest_log",
	ShowLine:      true,
	EncodeLevel:   "CapitalLevelEncoder",
	StacktraceKey: "stacktrace",
	LogInConsole:  false,
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"` //debug ,info,warn,error,panic,fatal
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	LinkName      string `mapstructure:"link-name" json:"link-name" yaml:"link-name"`
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"`
}

// IsExist config file is exist
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove remove config file
func Remove() error {
	return getViperConfig().Remove()
}

// Recover
func Recover() error {
	b, err := json.Marshal(CONFIG)
	if err != nil {
		return err
	}
	return getViperConfig().Recover(b)
}

// getViperConfig get viper config
func getViperConfig() viper_server.ViperConfig {
	configName := "zap"

	showLine := strconv.FormatBool(CONFIG.ShowLine)
	logInConsole := strconv.FormatBool(CONFIG.LogInConsole)

	config := viper_server.ViperConfig{
		Debug:     true,
		Name:      configName,
		Directory: consts.ConfDir,
		Type:      consts.ConfigType,

		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			vi.SetConfigName(configName)
			return nil
		},

		Default: []byte(`
{
	"level": "` + CONFIG.Level + `",
	"format": "` + CONFIG.Format + `",
	"prefix": "` + CONFIG.Prefix + `",
	"director": "` + CONFIG.Director + `",
	"link-name": "` + CONFIG.LinkName + `",
	"show-line": ` + showLine + `,
	"encode-level": "` + CONFIG.EncodeLevel + `",
	"stacktrace-key": "` + CONFIG.StacktraceKey + `",
	"log-in-console": ` + logInConsole + `
}`),
	}

	return config
}
