package domain

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/snowlyg/helper/str"
	"regexp"
)

type BaseUser struct {
	Name     string `gorm:"index;not null; type:varchar(60)" json:"name"`
	Username string `gorm:"uniqueIndex;not null;type:varchar(60)" json:"username" validate:"required"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"intro"`
	Avatar   string `gorm:"type:varchar(1024)" json:"avatar"`
}

type Avatar struct {
	Avatar string `json:"avatar"`
}

type UserDetail struct {
	_domain.Model
	BaseUser

	SysRoles     []string                 `gorm:"-" json:"sysRoles"`
	ProjectRoles map[uint]consts.RoleType `gorm:"-" json:"projectRoles"`

	RoleIds []uint `gorm:"-" json:"role_ids"`
}

func (res *UserDetail) ToString() {
	if res.Avatar == "" {
		return
	}
	re := regexp.MustCompile("^http")
	if !re.MatchString(res.Avatar) {
		res.Avatar = str.Join("http://127.0.0.1:8085/upload/", res.Avatar)
	}
}
