package web

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	"strings"
)

// Initialize  initialize
func Initialize() error {
	var cover string
	if config.IsExist() {
		fmt.Println("Your web config is initialized , reinitialized web will cover your web config.")
		fmt.Println("Did you want to do it ?  [Y/N]")
		fmt.Scanln(&cover)
		switch strings.ToUpper(cover) {
		case "Y":
		case "N":
			return nil
		default:
		}
	}

	err := config.Remove()
	if err != nil {
		return err
	}

	err = initConfig()
	if err != nil {
		return err
	}
	fmt.Println("web iris-admin initialized finished!")
	return nil
}

func initConfig() error {
	var dbType string
	fmt.Println("Please choose your database type: ")
	fmt.Println("1. mysql (only support mysql now)")
	fmt.Scanln(&dbType)
	switch dbType {
	case "1":
		config.CONFIG.System.DbType = "mysql"
		if err := database.Init(); err != nil {
			return err
		}
	default:
		config.CONFIG.System.DbType = "mysql"
		if err := database.Init(); err != nil {
			return err
		}
	}

	var systemTimeFormat, systemAddr string
	fmt.Println("Please input your system timeformat: ")
	fmt.Printf("System timeformat is '%s'\n", config.CONFIG.System.TimeFormat)
	fmt.Scanln(&systemTimeFormat)
	if systemTimeFormat != "" {
		config.CONFIG.System.TimeFormat = systemTimeFormat
	}

	fmt.Println("Please input your system addr: ")
	fmt.Printf("System addr is '%s'\n", config.CONFIG.System.Addr)
	fmt.Scanln(&systemAddr)
	if systemAddr != "" {
		config.CONFIG.System.Addr = systemAddr
	}
	err := viper_server.Init(config.GetViperConfig())
	if err != nil {
		return err
	}
	return nil
}
