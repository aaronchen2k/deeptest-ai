package operation

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func GetMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20211214120700_create_oplogs_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.SysOplog{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(model.SysOplog.TableName)
		},
	}
}
