package repo

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/casbin"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
	"gorm.io/gorm"
)

type RoleRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`
}

func (r *RoleRepo) Paginate(req v1.RolePageReq) (ret _domain.PageData, err error) {
	var count int64
	db := r.DB.Model(&model.SysRole{}).
		Where("NOT deleted")

	err = db.Count(&count).Error
	if err != nil {
		_logUtils.Errorf("count report error %s", err.Error())
		return
	}

	db.Scopes(r.PaginateScope(req.Page, req.PageSize, req.Order, req.Field))

	if req.Name != "" {
		db.Where("name = ?", req.Name)
	}

	results := make([]model.SysRole, 0)
	err = db.Find(&results).Error
	if err != nil {
		_logUtils.Errorf("query role error %s", err.Error())
		return
	}

	ret.Populate(results, count, req.Page, req.PageSize)

	return
}

func (r *RoleRepo) List() (pos []model.SysRole, err error) {
	err = r.DB.Model(&model.SysRole{}).
		Where("NOT deleted").
		Find(&pos).Error

	return
}

func (r *RoleRepo) Get(id uint) (po model.SysRole, err error) {
	err = r.DB.Where("id = ?", id).
		First(&po).Error

	return
}

func (r *RoleRepo) Create(po *model.SysRole) (id uint, err error) {
	err = r.DB.Create(po).Error

	return
}

func (r *RoleRepo) FindByName(name string) (ret *model.SysRole, err error) {
	err = r.DB.Model(&ret).
		Where("name = ?", name).
		First(&ret).
		Error

	return
}

func (r *RoleRepo) GetPermsForRole() (ret [][]string, err error) {
	perms := []model.SysPerm{}

	err = r.DB.
		Model(&model.SysPerm{}).
		Find(&perms).Error
	if err != nil {
		return
	}

	for _, perm := range perms {
		ret = append(ret, []string{perm.Name, perm.Act})
	}

	return
}

func (r *RoleRepo) AddPermForRole(name string, perms [][]string) (err error) {
	oldPerms := casbin.Instance().GetPermissionsForUser(name)

	_, err = casbin.Instance().RemovePolicies(oldPerms)
	if err != nil {
		_logUtils.Error(err.Error())
		return
	}

	if len(perms) == 0 {
		_logUtils.Debug("没有权限")
		return
	}

	var newPerms [][]string
	for _, perm := range perms {
		newPerms = append(newPerms, append([]string{name}, perm...))
	}

	_logUtils.Zap.Info("添加权限到角色", _str.Strings("新权限", newPerms))
	_, err = casbin.Instance().AddPolicies(newPerms)
	if err != nil {
		_logUtils.Error(err.Error())
		return err
	}

	return nil

	return
}

func (r *RoleRepo) GetRoleNames() (roleNames []string, err error) {
	err = r.DB.Model(&model.SysRole{}).
		Select("name").
		Find(&roleNames).Error
	if err != nil {
		_logUtils.Error(err.Error())
		return
	}

	return
}

func (r *RoleRepo) FindInName(names []string) (roles []v1.RoleResp, error error) {
	err := r.DB.Model(&model.SysRole{}).Where("name in (?)", names).Find(&roles).Error
	if err != nil {
		return
	}
	return
}
