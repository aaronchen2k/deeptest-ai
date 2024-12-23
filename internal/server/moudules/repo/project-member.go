package repo

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"gorm.io/gorm"
)

type ProjectMemberRoleRepo struct {
	DB        *gorm.DB `inject:""`
	*BaseRepo `inject:""`
}

func (r *ProjectMemberRoleRepo) GetProjectIdsByUser(userId uint) (projectIds []uint, err error) {
	r.DB.Model(&model.ProjectMember{}).
		Select("project_id").Where("user_id = ?", userId).Scan(&projectIds)

	return
}
