package repo

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
	"gorm.io/gorm"
	"strings"
)

type SetRepo struct {
	DB          *gorm.DB `inject:""`
	*BaseRepo   `inject:""`
	ProjectRepo *ProjectRepo `inject:""`
	UserRepo    *UserRepo    `inject:""`
}

func (r *SetRepo) Paginate(req v1.ReqPaginate, projectId int) (data _domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.TestSet{}).
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

	plans := make([]*model.TestSet, 0)

	err = db.
		Scopes(r.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&plans).Error
	if err != nil {
		return
	}

	data.Populate(plans, count, req.Page, req.PageSize)

	return
}

func (r *SetRepo) Get(id uint) (po model.TestSet, err error) {
	err = r.DB.Model(&model.TestSet{}).Where("id = ?", id).First(&po).Error
	if err != nil {
		return
	}

	return
}

func (r *SetRepo) Create(po model.TestSet) (ret model.TestSet, err error) {
	err = r.DB.Model(&model.TestSet{}).Create(&po).Error
	if err != nil {
		return
	}

	return
}

func (r *SetRepo) Update(req model.TestSet) (err error) {
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

func (r *SetRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.TestSet{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		return
	}

	return
}
