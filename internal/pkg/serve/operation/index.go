package operation

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
)

func init() {
	err := viper_server.Init(getViperConfig())
	if err != nil {
		panic(err)
	}
}

// CreateOplog
func CreateOplog(ol *model.SysOplog) error {
	err := database.GetInstance().Model(&model.SysOplog{}).Create(ol).Error
	if err != nil {
		return err
	}
	return nil
}
