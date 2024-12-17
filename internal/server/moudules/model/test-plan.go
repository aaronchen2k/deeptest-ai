package model

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"time"
)

type TestPlan struct {
	BaseModel

	Version float64 `json:"version" yaml:"version"`
	Name    string  `json:"name" yaml:"name"`
	Desc    string  `json:"desc" yaml:"desc"`

	Status   consts.TestStatus `json:"status"`
	ExecTime *time.Time        `gorm:"-" json:"execTime"`

	ProjectId uint `json:"projectId"`
}

func (TestPlan) TableName() string {
	return "tst_plans"
}
