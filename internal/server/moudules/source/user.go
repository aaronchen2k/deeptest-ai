package source

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gookit/color"
	"gorm.io/gorm"
)

var UserSrc = UserSource{}

type UserSource struct {
}

func GetUserMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20241206000000_create_sys_users_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.SysUser{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(model.SysUser.TableName)
		},
	}
}

func (s UserSource) Init() error {
	return database.GetInstance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&model.SysUser{}).Where("id IN ?", []int{1}).Find(&[]model.SysUser{}).RowsAffected == 1 {
			color.Danger.Println("\nusers 表的初始数据已存在!")
			return nil
		}
		sources, err := s.GetSources()
		if err != nil {
			return err
		}

		repo := repo.UserRepo{
			DB: database.GetInstance(),
		}

		for _, source := range sources {
			exist, _ := repo.FindByName(source.Name)
			if exist.ID > 0 {
				continue
			}

			_, err := repo.Create(source)
			if err != nil { // 遇到错误时回滚事务
				return err
			}
		}

		color.Info.Println("\n--> users 表初始数据成功!")

		return nil
	})
}

func (s UserSource) GetSources() ([]*v1.UserReq, error) {
	roleRepo := repo.RoleRepo{
		DB: database.GetInstance(),
	}

	roleNames, err := roleRepo.GetRoleNames()
	if err != nil {
		return []*v1.UserReq{}, err
	}

	var users []*v1.UserReq
	users = append(users, &v1.UserReq{
		BaseUser: domain.BaseUser{
			Name:     "管理员",
			Username: "admin",
			Intro:    "管理员",
			Avatar:   "/static/avatar.png",
		},
		Password:  "P2ssw0rd",
		RoleNames: roleNames,
	})

	return users, nil
}
