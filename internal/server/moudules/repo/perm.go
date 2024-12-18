package repo

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"gorm.io/gorm"
)

type PermRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`
}

func (r *PermRepo) Paginate(req v1.PermPageReq) (ret _domain.PageData, err error) {
	var count int64
	db := r.DB.Model(&model.SysPerm{}).
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
	if req.DisplayName != "" {
		db.Where("display_name = ?", req.DisplayName)
	}

	results := make([]model.SysPerm, 0)
	err = db.Find(&results).Error
	if err != nil {
		_logUtils.Errorf("query perm error %s", err.Error())
		return
	}

	ret.Populate(results, count, req.Page, req.PageSize)

	return
}

func (r *PermRepo) List() (pos []model.SysPerm, err error) {
	err = r.DB.Model(&model.SysPerm{}).
		Where("NOT deleted").
		Find(&pos).Error

	return
}

func (r *PermRepo) Get(id uint) (po model.SysPerm, err error) {
	err = r.DB.Where("id = ?", id).
		First(&po).Error

	return
}
