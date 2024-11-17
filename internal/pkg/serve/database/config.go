package database

import (
	"encoding/json"
	"fmt"
	config2 "github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/spf13/viper"
)

var CONFIG_MYSQL = Mysql{
	Path:         "127.0.0.1:3306",
	Config:       "charset=utf8mb4&parseTime=True&loc=Local",
	DbName:       "deeptest-db",
	Username:     "root",
	Password:     "",
	MaxIdleConns: 0,
	MaxOpenConns: 0,
	LogMode:      true,
	LogZap:       "zap",
}
var CONFIG_SQLITE = Sqlite{
	DbName:       "deeptest",
	MaxIdleConns: 0,
	MaxOpenConns: 0,
	LogMode:      true,
	LogZap:       "zap",
}

type Mysql struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Config       string `mapstructure:"conf" json:"conf" yaml:"conf"`
	DbName       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`
	LogZap       string `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"` //silent,error,warn,info,zap
}

// Dsn return mysql dsn
func (m *Mysql) Dsn() string {
	return fmt.Sprintf("%s%s?%s", m.BaseDsn(), m.DbName, m.Config)
}

// Dsn return
func (m *Mysql) BaseDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/", m.Username, m.Password, m.Path)
}

type Sqlite struct {
	DbName       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`
	LogZap       string `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"` //silent,error,warn,info,zap
}

// IsExist conf file is exist
func IsExist() bool {
	return GetViperConfig().IsFileExist()
}

// Remove remove conf file
func Remove() error {
	return GetViperConfig().Remove()
}

// Recover
func Recover() error {
	b, err := json.Marshal(CONFIG_MYSQL)
	if err != nil {
		return err
	}
	return GetViperConfig().Recover(b)
}

func GetViperConfig() viper_server.ViperConfig {
	if config2.CONFIG.System.DatabaseType == "sqlite" {
		return getViperConfigSqlite()
	} else {
		return getViperConfigMySql()
	}
}

// getViperConfigMySql get viper conf
func getViperConfigMySql() viper_server.ViperConfig {
	configName := "mysql"
	mxIdleConns := fmt.Sprintf("%d", CONFIG_MYSQL.MaxIdleConns)
	mxOpenConns := fmt.Sprintf("%d", CONFIG_MYSQL.MaxOpenConns)
	logMode := fmt.Sprintf("%t", CONFIG_MYSQL.LogMode)

	return viper_server.ViperConfig{
		Debug:     true,
		Directory: consts.ConfDir,
		Name:      configName,
		Type:      consts.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG_MYSQL); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			// watch conf file change
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("Config file changed:", e.Name)
			})
			vi.WatchConfig()
			return nil
		},
		//
		Default: []byte(`
{
	"path": "` + CONFIG_MYSQL.Path + `",
	"conf": "` + CONFIG_MYSQL.Config + `",
	"db-name": "` + CONFIG_MYSQL.DbName + `",
	"username": "` + CONFIG_MYSQL.Username + `",
	"password": "` + CONFIG_MYSQL.Password + `",
	"max-idle-conns": ` + mxIdleConns + `,
	"max-open-conns": ` + mxOpenConns + `,
	"log-mode": ` + logMode + `,
	"log-zap": "` + CONFIG_MYSQL.LogZap + `"
}`),
	}
}

// getViperConfigSqlite get viper conf
func getViperConfigSqlite() viper_server.ViperConfig {
	configName := "sqlite"
	mxIdleConns := fmt.Sprintf("%d", CONFIG_SQLITE.MaxIdleConns)
	mxOpenConns := fmt.Sprintf("%d", CONFIG_SQLITE.MaxOpenConns)
	logMode := fmt.Sprintf("%t", CONFIG_SQLITE.LogMode)

	return viper_server.ViperConfig{
		Debug:     true,
		Directory: consts.ConfDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG_SQLITE); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			// watch conf file change
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("Config file changed:", e.Name)
			})
			vi.WatchConfig()
			return nil
		},
		//
		Default: []byte(`
{
	"db-name": "` + CONFIG_SQLITE.DbName + `",
	"max-idle-conns": ` + mxIdleConns + `,
	"max-open-conns": ` + mxOpenConns + `,
	"log-mode": ` + logMode + `,
	"log-zap": "` + CONFIG_SQLITE.LogZap + `"
}`),
	}
}
