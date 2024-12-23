package service

import (
	"errors"
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/auth"
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type AccountService struct {
	AccountRepo *repo.AccountRepo `inject:""`
	UserRepo    *repo.UserRepo    `inject:""`
	ProfileRepo *repo.ProfileRepo `inject:""`
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func (s *AccountService) GetAccessToken(req *v1.LoginReq) (token string, id uint, err error) {
	admin, err := s.UserRepo.GetPasswordByUserName(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if admin == nil || admin.ID == 0 {
		err = consts.ErrUserNameOrPassword
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		_logs.Zap.Error("用户名或密码错误", zap.String("密码:", req.Password), zap.String("bcrypt.CompareHashAndPassword()", err.Error()))

		err = consts.ErrUserNameOrPassword
		return
	}

	expiresAt := time.Now().Local().Add(time.Duration(config.CONFIG.SessionTimeout) * time.Minute).Unix()
	claims := auth.New(&auth.Multi{
		Id:            admin.ID,
		Username:      req.Username,
		AuthorityIds:  admin.AuthorityIds,
		AuthorityType: auth.AdminAuthority,
		LoginType:     auth.LoginTypeWeb,
		AuthType:      auth.AuthPwd,
		ExpiresAt:     expiresAt,
	})

	token, _, err = auth.AuthDriver.GenerateToken(claims)
	if err != nil {
		_logs.Zap.Error(err.Error())
		return
	}

	id = admin.ID

	return
}

func (s *AccountService) DeleteToken(token string) (err error) {
	err = auth.AuthDriver.DelUserTokenCache(token)

	if err != nil {
		_logs.Error(err.Error())
		return
	}

	return
}

func (s *AccountService) CleanToken(authorityType int, userId string) (err error) {
	err = auth.AuthDriver.CleanUserTokenCache(authorityType, userId)
	if err != nil {
		_logs.Error(err.Error())
		return
	}

	return
}

func (s *AccountService) GetInfo(userId uint) (info domain.UserDetail, err error) {
	info, err = s.AccountRepo.GetInfo(userId)

	return
}

func (s *AccountService) GetCodes(userId uint) (ret []string, err error) {
	return
}

func (s *AccountService) UpdateUserProject(req *v1.UpdateUserProjectReq, userId uint) (err error) {
	err = s.ProfileRepo.UpdateUserProject(req.ProjectId, userId)

	return
}
