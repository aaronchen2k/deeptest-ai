package config

import (
	"encoding/json"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

var CONFIG = Web{
	FileMaxSize:    1024, // upload file size limit 1024M
	SessionTimeout: 60,   // session timeout after 60 Minute
	Cors: Cors{
		AccessOrigin:        "*",
		AccessHeaders:       "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id",
		AccessMethods:       "POST,GET,OPTIONS,DELETE,PUT",
		AccessExposeHeaders: "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type",
		AccessCredentials:   "true",
	},
	Except: Route{
		Uri:    "",
		Method: "",
	},
	System: System{
		Tls:        false,
		Level:      "debug",
		Addr:       "0.0.0.0:9085",
		DbType:     "mysql",
		TimeFormat: "2006-01-02 15:04:05",

		ChatchatUrl: "http://127.0.0.1:7861/",
	},
	Limit: Limit{
		Disable: true,
		Limit:   0,
		Burst:   5,
	},
	Captcha: Captcha{
		KeyLong:   4,
		ImgWidth:  240,
		ImgHeight: 80,
	},
}

type Web struct {
	FileMaxSize    int64   `mapstructure:"file-max-size" json:"file-max-size" yaml:"file-max-siz"`
	SessionTimeout int64   `mapstructure:"session-timeout" json:"session-timeout" yaml:"session-timeout"`
	Except         Route   `mapstructure:"except" json:"except" yaml:"except"`
	System         System  `mapstructure:"system" json:"system" yaml:"system"`
	Limit          Limit   `mapstructure:"limit" json:"limit" yaml:"limit"`
	Captcha        Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Cors           Cors    `mapstructure:"cors" json:"cors" yaml:"cors"`
}

type Cors struct {
	AccessOrigin        string `mapstructure:"access-origin" json:"burst" access-origin:"access-origin"`
	AccessHeaders       string `mapstructure:"access-headers" json:"access-headers" yaml:"access-headers"`
	AccessMethods       string `mapstructure:"access-methods" json:"access-methods" yaml:"access-methods"`
	AccessExposeHeaders string `mapstructure:"access-expose-headers" json:"access-expose-headers" yaml:"access-expose-headers"`
	AccessCredentials   string `mapstructure:"access-credentials" json:"access-credentials" yaml:"access-credentials"`
}
type Route struct {
	Uri    string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Method string `mapstructure:"method" json:"method" yaml:"method"`
}

type Captcha struct {
	KeyLong   int `mapstructure:"key-long" json:"key-long" yaml:"key-long"`
	ImgWidth  int `mapstructure:"img-width" json:"img-width" yaml:"img-width"`
	ImgHeight int `mapstructure:"img-height" json:"img-height" yaml:"img-height"`
}

type Limit struct {
	Disable bool    `mapstructure:"disable" json:"disable" yaml:"disable"`
	Limit   float64 `mapstructure:"limit" json:"limit" yaml:"limit"`
	Burst   int     `mapstructure:"burst" json:"burst" yaml:"burst"`
}

type System struct {
	Tls          bool   `mapstructure:"tls" json:"tls" yaml:"tls"`       // debug,release,test
	Level        string `mapstructure:"level" json:"level" yaml:"level"` // debug,release,test
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`
	StaticPrefix string `mapstructure:"static-prefix" json:"static-prefix" yaml:"static-prefix"`
	WebPrefix    string `mapstructure:"web-prefix" json:"web-prefix" yaml:"web-prefix"`
	DbType       string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`
	TimeFormat   string `mapstructure:"time-format" json:"time-format" yaml:"time-format"`

	ChatchatUrl string `mapstructure:"chatchat-url" json:"chatchat-url" yaml:"chatchat-url"`
	LLmUrl      string `mapstructure:"llm-url" json:"llm-url" yaml:"llm-url"`

	DatabaseType string `mapstructure:"database-type" json:"database-type" yaml:"database-type"`
}

// SetDefaultAddrAndTimeFormat
func SetDefaultAddrAndTimeFormat() {
	if CONFIG.System.Addr == "" {
		CONFIG.System.Addr = "0.0.0.0:9085"
	}

	if CONFIG.System.TimeFormat == "" {
		CONFIG.System.TimeFormat = "2006-01-02 15:04:05"
	}
}

// ToStaticUrl
func ToStaticUrl(uri string) string {
	path := filepath.Join(CONFIG.System.Addr, CONFIG.System.StaticPrefix, uri)
	if CONFIG.System.Tls {
		return filepath.ToSlash(_str.Join("https://", path))
	}
	return filepath.ToSlash(_str.Join("http://", path))
}

// IsExist config file is exist
func IsExist() bool {
	return GetViperConfig().IsFileExist()
}

// Remove remove config file
func Remove() error {
	return GetViperConfig().Remove()
}

// Recover
func Recover() error {
	b, err := json.Marshal(CONFIG)
	if err != nil {
		return err
	}
	return GetViperConfig().Recover(b)
}

// GetViperConfig get viper config
func GetViperConfig() viper_server.ViperConfig {
	maxSize := strconv.FormatInt(CONFIG.FileMaxSize, 10)
	sessionTimeout := strconv.FormatInt(CONFIG.SessionTimeout, 10)
	keyLong := strconv.FormatInt(int64(CONFIG.Captcha.KeyLong), 10)
	imgWidth := strconv.FormatInt(int64(CONFIG.Captcha.ImgWidth), 10)
	imgHeight := strconv.FormatInt(int64(CONFIG.Captcha.ImgHeight), 10)
	limit := strconv.FormatInt(int64(CONFIG.Limit.Limit), 10)
	burst := strconv.FormatInt(int64(CONFIG.Limit.Burst), 10)
	disable := strconv.FormatBool(CONFIG.Limit.Disable)
	tls := strconv.FormatBool(CONFIG.System.Tls)
	configName := "web"
	return viper_server.ViperConfig{
		Debug:     true,
		Directory: consts.ConfDir,
		Name:      configName,
		Type:      consts.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			// watch config file change
			vi.SetConfigName(configName)
			return nil
		},
		//
		Default: []byte(`
{
	"file-max-size": ` + maxSize + `,
	"session-timeout": ` + sessionTimeout + `,
	"except":
		{ 
			"uri": "` + CONFIG.Except.Uri + `",
			"method": "` + CONFIG.Except.Method + `"
		},
	"cors":
		{ 
			"access-origin": "` + CONFIG.Cors.AccessOrigin + `",
			"access-headers": "` + CONFIG.Cors.AccessHeaders + `",
			"access-methods": "` + CONFIG.Cors.AccessMethods + `",
			"access-expose-headers": "` + CONFIG.Cors.AccessExposeHeaders + `",
			"access-credentials": "` + CONFIG.Cors.AccessCredentials + `"
		},
	"captcha":
		{
		"key-long": ` + keyLong + `,
		"img-width": ` + imgWidth + `,
		"img-height": ` + imgHeight + `
		},
	"limit":
		{
			"limit": ` + limit + `,
			"disable": ` + disable + `,
			"burst": ` + burst + `
		},
	"system":
		{
			"tls": ` + tls + `,
			"level": "` + CONFIG.System.Level + `",
			"addr": "` + CONFIG.System.Addr + `",
			"db-type": "` + CONFIG.System.DbType + `",
			"time-format": "` + CONFIG.System.TimeFormat + `",
    		"chatchat-url": "` + CONFIG.System.ChatchatUrl + `"
		}
 }`),
	}
}
