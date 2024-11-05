package repo

import (
	"errors"
	"fmt"
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/casbin"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

type UserRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`
}

func (r *UserRepo) Paginate(req v1.UserPageReq) (ret _domain.PageData, err error) {
	var count int64
	db := r.DB.Model(&model.SysUser{}).
		Where("NOT deleted")

	err = db.Count(&count).Error
	if err != nil {
		_logUtils.Errorf("count report error %s", err.Error())
		return
	}

	db.Scopes(PaginateScope(req.Page, req.PageSize, req.Order, req.Field))

	if req.Name != "" {
		db.Where("name = ?", req.Name)
	}

	results := make([]model.SysUser, 0)
	err = db.Find(&results).Error
	if err != nil {
		_logUtils.Errorf("query user error %s", err.Error())
		return
	}

	ret.Populate(results, count, req.Page, req.PageSize)

	return
}

func (r *UserRepo) List() (pos []model.SysUser, err error) {
	err = r.DB.Model(&model.SysUser{}).
		Where("NOT deleted").
		Find(&pos).Error

	return
}

func (r *UserRepo) Get(id uint) (po model.SysUser, err error) {
	err = r.DB.Where("id = ?", id).
		First(&po).Error

	return
}

func (r *UserRepo) GetPasswordByUserName(username string) (
	ret *v1.LoginResp, err error) {

	user := &v1.LoginResp{}

	err = r.DB.Model(&model.SysUser{}).
		Select("id, password").
		Where("username = ?", username).
		First(user).Error
	if err != nil {
		return
	}

	user.AuthorityIds, err = r.GetUserRoleNames(fmt.Sprintf("%v", user.ID))
	if err != nil {
		return
	}

	ret = user

	return
}

func (r *UserRepo) GetUserRoleNames(userId string) (roleNames []string, err error) {
	roleNames, err = casbin.Instance().GetRolesForUser(userId)
	if err != nil {
		return
	}

	return
}

func (r *UserRepo) Create(req *v1.UserReq) (user model.SysUser, err error) {
	_, err = r.FindByName(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		err = consts.ErrUserNameInvalid
		return
	}

	user = model.SysUser{BaseUser: req.BaseUser, RoleNames: req.RoleNames}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		_logUtils.Error(err.Error())
		return
	}

	_logUtils.Zap.Info("添加用户", zap.String("hash:", req.Password), zap.ByteString("hash:", hash))

	user.Password = string(hash)

	err = r.DB.Create(&user).Error
	if err != nil {
		_logUtils.Error(err.Error())
		return
	}

	err = r.AddRoleForUser(user)
	if err != nil {
		return
	}

	return
}

func (r *UserRepo) Update(user model.SysUser) (err error) {
	err = r.DB.Save(&user).Error
	if err != nil {
		_logUtils.Error(err.Error())
		return
	}

	return
}

func (r *UserRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.SysUser{}).
		Where("id = ?", id).
		Update("deleted", true).Error

	return
}

func (r *UserRepo) AddRoleForUser(user model.SysUser) (err error) {
	userId := strconv.FormatUint(uint64(user.ID), 10)
	oldRoleNames, err := r.GetUserRoleNames(userId)
	if err != nil {
		return err
	}

	if len(oldRoleNames) > 0 {
		_, err = casbin.Instance().DeleteRolesForUser(userId)
		if err != nil {
			_logUtils.Error(err.Error())
			return err
		}
	}

	if len(user.RoleNames) == 0 {
		return nil
	}

	var roleNames []string
	roleNames = append(roleNames, user.RoleNames...)

	_, err = casbin.Instance().AddRolesForUser(userId, roleNames)
	if err != nil {
		_logUtils.Error(err.Error())
		return err
	}

	return
}

func (r *UserRepo) FindByName(name string) (ret *model.SysUser, err error) {
	err = r.DB.Model(&ret).
		Where("name = ?", name).
		First(&ret).
		Error

	return
}
