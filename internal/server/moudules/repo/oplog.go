package repo

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"gorm.io/gorm"
)

type OplogRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`
}

func (r *OplogRepo) Paginate(req v1.OptLogPageReq) (ret _domain.PageData, err error) {
	var count int64
	db := r.DB.Model(&model.SysOplog{}).
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
	if req.UserName != "" {
		db.Where("user_name = ?", req.UserName)
	}

	results := make([]model.SysOplog, 0)
	err = db.Find(&results).Error
	if err != nil {
		_logUtils.Errorf("query oplog error %s", err.Error())
		return
	}

	ret.Populate(results, count, req.Page, req.PageSize)

	return
}

func (r *OplogRepo) List() (pos []model.SysOplog, err error) {
	err = r.DB.Model(&model.SysOplog{}).
		Where("NOT deleted").
		Find(&pos).Error

	return
}
