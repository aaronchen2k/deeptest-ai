package model

import "github.com/deeptest-com/deeptest-next/internal/pkg/consts"

type ProjectRole struct {
	BaseModel

	Name        consts.RoleType `gorm:"uniqueIndex;not null; type:varchar(256)" json:"name"`
	DisplayName string          `gorm:"type:varchar(256)" json:"displayName" comment:"显示名称"`
	Description string          `gorm:"type:varchar(256)" json:"description" comment:"描述"`
}

func (ProjectRole) TableName() string {
	return "biz_project_role"
}
