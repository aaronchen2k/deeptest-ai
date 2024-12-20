package model

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
)

type Project struct {
	BaseModel
	domain.ProjectBase

	UpdatedUser string `json:"updatedUser" gorm:"omitempty"`
}

func (Project) TableName() string {
	return "biz_project"
}
