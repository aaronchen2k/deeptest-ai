package repo

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
	"gorm.io/gorm"
	"strings"
)

type CaseRepo struct {
	DB          *gorm.DB `inject:""`
	*BaseRepo   `inject:""`
	ProjectRepo *ProjectRepo `inject:""`
	UserRepo    *UserRepo    `inject:""`
}

func (r *CaseRepo) Paginate(req v1.ReqPaginate, projectId int) (data _domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.TestCase{}).
		Where("project_id = ? AND NOT deleted",
			projectId)

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.Enabled != "" {
		db = db.Where("disabled = ?", r.IsDisable(req.Enabled))
	}
	if req.Status != "" {
		db = db.Where("status IN (?)", strings.Split(req.Status, ","))
	}

	err = db.Count(&count).Error
	if err != nil {
		return
	}

	plans := make([]*model.TestCase, 0)

	err = db.
		Scopes(r.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&plans).Error
	if err != nil {
		return
	}

	data.Populate(plans, count, req.Page, req.PageSize)

	return
}

func (r *CaseRepo) Get(id uint) (scenario model.TestCase, err error) {
	err = r.DB.Model(&model.TestCase{}).Where("id = ?", id).First(&scenario).Error
	if err != nil {
		return
	}

	return
}

func (r *CaseRepo) Create(scenario model.TestCase) (ret model.TestCase, bizErr *_domain.BizErr) {
	err := r.DB.Model(&model.TestCase{}).Create(&scenario).Error
	if err != nil {
		bizErr = &_domain.BizErr{Code: _domain.SystemErr.Code}

		return
	}

	return
}

func (r *CaseRepo) Update(req model.TestCase) (err error) {
	values := map[string]interface{}{
		"name":   req.Name,
		"desc":   req.Desc,
		"status": req.Status,

		"disabled": req.Disabled,
		"deleted":  req.Deleted,
	}
	err = r.DB.Model(&req).Where("id = ?", req.ID).Updates(values).Error
	if err != nil {
		return
	}

	return
}

func (r *CaseRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.TestCase{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		return
	}

	return
}
