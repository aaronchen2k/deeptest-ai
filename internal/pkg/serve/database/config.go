package database

import (
	"encoding/json"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/spf13/viper"
	"strconv"
)

var CONFIG_POSTGRES = Postgres{
	Host:         "127.0.0.1",
	Port:         5432,
	DbName:       "deeptest",
	Username:     "root",
	Password:     "",
	MaxIdleConns: 0,
	MaxOpenConns: 0,
	LogMode:      true,
	LogZap:       "zap",
}
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

type Postgres struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	Config       string `mapstructure:"conf" json:"conf" yaml:"conf"`
	DbName       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`
	LogZap       string `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"` //silent,error,warn,info,zap
}

func (m *Postgres) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		m.Host, m.Port, m.Username, m.Password, m.DbName)
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

func (m *Mysql) Dsn() string {
	return fmt.Sprintf("%s%s?%s", m.BaseDsn(), m.DbName, m.Config)
}
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
	if config.CONFIG.System.DatabaseType == "postgres" {
		return getViperConfigPostgres()
	} else if config.CONFIG.System.DatabaseType == "sqlite" {
		return getViperConfigSqlite()
	} else {
		return getViperConfigMySql()
	}
}

func getViperConfigPostgres() viper_server.ViperConfig {
	configName := "postgres"
	mxIdleConns := fmt.Sprintf("%d", CONFIG_POSTGRES.MaxIdleConns)
	mxOpenConns := fmt.Sprintf("%d", CONFIG_POSTGRES.MaxOpenConns)
	logMode := fmt.Sprintf("%t", CONFIG_POSTGRES.LogMode)

	conf := viper_server.ViperConfig{
		Debug:     true,
		Directory: consts.ConfDir,
		Name:      configName,
		Type:      consts.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG_POSTGRES); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}

			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("Config file changed:", e.Name)
			})
			vi.WatchConfig()
			return nil
		},

		Default: []byte(`
{
	"host": "` + CONFIG_POSTGRES.Host + `",
	"port": ` + strconv.Itoa(CONFIG_POSTGRES.Port) + `,
	"conf": "` + CONFIG_POSTGRES.Config + `",
	"db-name": "` + CONFIG_POSTGRES.DbName + `",
	"username": "` + CONFIG_POSTGRES.Username + `",
	"password": "` + CONFIG_POSTGRES.Password + `",
	"max-idle-conns": ` + mxIdleConns + `,
	"max-open-conns": ` + mxOpenConns + `,
	"log-mode": ` + logMode + `,
	"log-zap": "` + CONFIG_POSTGRES.LogZap + `"
}`),
	}

	return conf
}

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
