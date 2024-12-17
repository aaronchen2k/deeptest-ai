package repo

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
	"gorm.io/gorm"
	"strings"
)

type PlanRepo struct {
	DB          *gorm.DB `inject:""`
	*BaseRepo   `inject:""`
	ProjectRepo *ProjectRepo `inject:""`
	UserRepo    *UserRepo    `inject:""`
}

func (r *PlanRepo) Paginate(req v1.ReqPaginate, projectId int) (data _domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.TestPlan{}).
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

	plans := make([]*model.TestPlan, 0)

	err = db.
		Scopes(r.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&plans).Error
	if err != nil {
		return
	}

	data.Populate(plans, count, req.Page, req.PageSize)

	return
}

func (r *PlanRepo) Get(id uint) (scenario model.TestPlan, err error) {
	err = r.DB.Model(&model.TestPlan{}).Where("id = ?", id).First(&scenario).Error
	if err != nil {
		return scenario, err
	}

	return scenario, nil
}

func (r *PlanRepo) FindByName(name string, id uint) (scenario model.TestPlan, err error) {
	db := r.DB.Model(&model.TestPlan{}).
		Where("name = ? AND NOT deleted", name)

	if id > 0 {
		db.Where("id != ?", id)
	}

	db.First(&scenario)

	return
}

func (r *PlanRepo) Create(scenario model.TestPlan) (ret model.TestPlan, err error) {
	err = r.DB.Model(&model.TestPlan{}).Create(&scenario).Error
	if err != nil {
		return
	}

	return
}

func (r *PlanRepo) Update(req model.TestPlan) (err error) {
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

func (r *PlanRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.TestPlan{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		return
	}

	return
}
