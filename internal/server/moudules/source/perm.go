package source

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gookit/color"
	"gorm.io/gorm"
)

type PermSource struct {
	PermRepo *repo.PermRepo `inject:""`
	routes   []map[string]string
}

func NewPermSource(routes []map[string]string) *PermSource {
	return &PermSource{
		routes: routes,
	}
}

func GetPermMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240606000000_create_sys_permissions_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.SysPerm{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(model.SysPerm.TableName)
		},
	}
}

func (s *PermSource) Init() error {
	if s.getSources() == nil {
		return nil
	}

	return database.GetInstance().Transaction(func(tx *gorm.DB) error {
		err := tx.Unscoped().Where("1 = 1").Delete(&model.SysPerm{}).Error
		if err != nil {
			return err
		}

		err = s.createBatch(tx, s.getSources())
		if err != nil { // 遇到错误时回滚事务
			return err
		}

		color.Info.Println("\n[Mysql] --> permissions 表初始数据成功!")

		return nil
	})
}

func (s *PermSource) getSources() model.PermCollection {
	perms := make(model.PermCollection, 0, len(s.routes))

	for _, permRoute := range s.routes {
		perm := model.SysPerm{
			BasePermission: domain.BasePermission{
				Name:        permRoute["path"],
				DisplayName: permRoute["name"],
				Description: permRoute["name"],
				Act:         permRoute["act"],
			}}
		perms = append(perms, perm)
	}

	return perms
}

func (s *PermSource) createBatch(tx *gorm.DB, perms model.PermCollection) (err error) {
	err = tx.Model(&model.SysPerm{}).
		CreateInBatches(&perms, 500).
		Error

	if err != nil {
		return
	}

	return
}
