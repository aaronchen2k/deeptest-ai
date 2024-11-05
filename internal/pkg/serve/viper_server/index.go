package viper_server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	ErrEmptyName = errors.New("config'name can't be empty value")
)

type ViperConfig struct {
	Debug     bool
	Directory string
	Name      string
	Type      string
	Default   []byte
	Watch     func(*viper.Viper) error
}

// getConfigFilePath
func (vc ViperConfig) getConfigFilePath() string {
	return filepath.Join(consts.ExecDir, vc.Directory, _str.Join(vc.Name, ".", vc.Type))
}

// getConfigFileDir
func (vc ViperConfig) getConfigFileDir() string {
	if vc.Directory == "" {
		return consts.ConfDir
	}
	return vc.Directory
}

// IsFileExist
func (vc ViperConfig) IsFileExist() bool {
	return _file.IsExist(vc.getConfigFilePath())
}

// Remove remove config file
func (vc ViperConfig) Remove() error {
	return _file.Remove(vc.getConfigFilePath())
}

// Recover recover config file content
func (vc ViperConfig) Recover(b []byte) error {
	_, err := _file.WriteBytes(vc.getConfigFilePath(), b)
	return err
}

// Init
func Init(vc ViperConfig) error {
	if vc.Name == "" {
		return ErrEmptyName
	}

	if vc.Type == "" {
		vc.Type = "yaml"
	}

	vc.Directory = vc.getConfigFileDir()

	filePath := vc.getConfigFilePath()
	if vc.Debug {
		fmt.Printf("\nthis config file's path is [%s]\n", filePath)
	}

	vi := viper.New()
	if vc.Debug {
		fmt.Printf("this config file's type is [%s]\n", vc.Type)
	}
	vi.SetConfigName(vc.Name)
	vi.SetConfigType(vc.Type)
	vi.AddConfigPath(vc.Directory)

	isExist := _file.IsExist(filePath)
	if !isExist {
		if vc.Debug {
			fmt.Printf("this config [%s] is not exist\n", filePath)
		}
		if vc.Directory != "./" {
			err := _file.InsureDir(filepath.Dir(filePath))
			if err != nil {
				return fmt.Errorf("create dir %s fail : %v", filePath, err)
			}
		}

		// ReadConfig
		if err := vi.ReadConfig(bytes.NewBuffer(vc.Default)); err != nil {
			if vc.Debug {
				fmt.Println(string(vc.Default))
			}
			return fmt.Errorf("read default config fail : %w ", err)
		}

		// WriteConfigAs
		if err := vi.WriteConfigAs(filePath); err != nil {
			return fmt.Errorf("write config to path fail: %w ", err)
		}

	} else {
		if vc.Debug {
			fmt.Printf("this config file [%s] is existed\n", filePath)
		}
		vi.SetConfigFile(filePath)
		err := vi.ReadInConfig()
		if err != nil {
			return fmt.Errorf("read config fail: %w ", err)
		}
	}

	err := vc.Watch(vi)
	if err != nil {
		return err
	}

	return nil
}
