package model

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
)

type SysRole struct {
	BaseModel
	domain.BaseRole

	Perms [][]string `gorm:"-" json:"perms"`
}

func (SysRole) TableName() string {
	return "sys_roles"
}
