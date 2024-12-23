package repo

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"gorm.io/gorm"
)

type ProfileRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB  `inject:""`
	RoleRepo  *RoleRepo `inject:""`
}

func (r *ProfileRepo) Get(userId uint) (profile model.SysUserProfile, err error) {
	db := r.DB.Model(&model.SysUserProfile{}).Where("user_id = ?", userId)
	err = db.First(&profile).Error
	return
}

func (r *ProfileRepo) UpdateUserProject(projectId, userId uint) (err error) {
	err = r.DB.Model(&model.SysUserProfile{}).Where("user_id = ?", userId).
		Updates(map[string]interface{}{"curr_project_id": projectId}).Error

	return
}
