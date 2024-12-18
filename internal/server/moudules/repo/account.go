package repo

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AccountRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`

	UserRepo *UserRepo `inject:""`
}

func (r AccountRepo) GetInfo(userId uint) (info domain.UserDetail, err error) {
	user, err := r.UserRepo.Get(userId)
	if err != nil {
		return
	}

	copier.CopyWithOption(&info, user, copier.Option{DeepCopy: true})

	r.UserRepo.GetSysRoles(&info)
	r.UserRepo.GetProjectRoles(&info)

	return
}
