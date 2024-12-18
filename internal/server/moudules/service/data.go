package service

import (
	"errors"
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	modules "github.com/deeptest-com/deeptest-next/internal/pkg/core/module"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/source"
	_logs "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrViperEmpty = errors.New("配置服务未初始化")
)

type DataService struct {
	DataRepo   *repo.DataRepo     `inject:""`
	UserRepo   *repo.UserRepo     `inject:""`
	UserSource *source.UserSource `inject:""`
	RoleSource *source.RoleSource `inject:""`
	PermSource *source.PermSource `inject:""`
}

// InitDB 创建数据库并初始化
func (s *DataService) InitDB(req v1.DataReq) error {
	err := s.DataRepo.DB.AutoMigrate(model.Models...)
	if err != nil {
		_logs.Errorf("迁移数据表错误", zap.String("错误:", err.Error()))
		return err
	}

	if req.ClearData {
		err := s.initData(
			s.PermSource,
			s.RoleSource,
			s.UserSource,
		)
		if err != nil {
			_logs.Errorf("填充数据错误", zap.String("错误:", err.Error()))
			return err
		}
	}

	if req.Sys.AdminPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Sys.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			_logs.Errorf("密码加密错误", zap.String("错误:", err.Error()))
			return nil
		}

		req.Sys.AdminPassword = string(hash)
		s.UserRepo.UpdatePasswordByName(consts.AdminUserName, req.Sys.AdminPassword)
	}

	return nil
}

// initDB 初始化数据
func (s *DataService) initData(InitDBFunctions ...modules.InitDBFunc) error {
	for _, v := range InitDBFunctions {
		err := v.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
