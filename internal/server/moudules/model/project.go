package model

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
)

type Project struct {
	BaseModel
	domain.ProjectBase

	AdminId   uint   `json:"adminId"`
	AdminName string `json:"adminName"`
}

func (Project) TableName() string {
	return "biz_project"
}
