package model

import "github.com/deeptest-com/deeptest-next/internal/pkg/consts"

type TestCase struct {
	BaseModel

	Version float64 `json:"version" yaml:"version"`
	Name    string  `json:"name" yaml:"name"`
	Desc    string  `json:"desc" yaml:"desc"`

	Status    consts.TestStatus `json:"status"`
	ProjectId uint              `json:"projectId"`
}

func (TestCase) TableName() string {
	return "tst_cases"
}
