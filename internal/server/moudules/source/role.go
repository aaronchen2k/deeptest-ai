package source

import (
	"errors"
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gookit/color"
	"gorm.io/gorm"
)

var RoleSrc = RoleSource{}

type RoleSource struct {
}

func GetRoleMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20241206000000_create_sys_roles_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.SysRole{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(model.SysRole.TableName)
		},
	}
}

func (s RoleSource) Init() error {
	db := database.GetInstance()
	//if db.Model(&model.SysRole{}).
	//	Where("id IN ?", []int{1}).
	//	Find(&[]model.SysRole{}).RowsAffected == 1 {
	//	color.Danger.Println("\n--> roles 表的初始数据已存在!")
	//	return nil
	//}

	sources, err := s.getSourcesWithPerms(db)
	if err != nil {
		return err
	}

	for _, source := range sources {
		_, err := s.Create(db, source)

		if err != nil { // 遇到错误时回滚事务
			return err
		}
	}

	color.Info.Println("\n--> roles 表初始数据成功!")
	return nil
}

func (s RoleSource) getSourcesWithPerms(tx *gorm.DB) ([]*v1.RoleReq, error) {
	repo := repo.RoleRepo{
		DB: tx,
	}

	perms, err := repo.GetPermsForRole()
	if err != nil {
		return []*v1.RoleReq{}, err
	}

	var sources []*v1.RoleReq
	sources = append(sources, &v1.RoleReq{
		BaseRole: domain.BaseRole{
			Name:        "admin",
			DisplayName: "管理员",
		},
		Perms: perms,
	})

	return sources, err
}

func (s RoleSource) Create(tx *gorm.DB, req *v1.RoleReq) (id uint, err error) {
	repo := repo.RoleRepo{DB: tx}

	_, err = repo.FindByName(req.Name)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			role := &model.SysRole{BaseRole: req.BaseRole}
			id, err = repo.Create(role)
			if err != nil {
				return
			}
		} else {
			err = consts.ErrRoleNameInvalid
			return
		}
	}

	err = repo.AddPermForRole(req.Name, req.Perms)
	if err != nil {
		return
	}

	return id, nil
}
