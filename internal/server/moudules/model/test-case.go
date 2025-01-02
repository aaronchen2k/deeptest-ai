package model

import "github.com/deeptest-com/deeptest-next/internal/pkg/consts"

type TestCase struct {
	BaseModel

	Version float64             `json:"version" yaml:"version"`
	Title   string              `json:"title" yaml:"title"`
	Type    consts.TreeNodeType `json:"type"`
	Desc    string              `json:"desc" yaml:"desc"`

	Status    consts.TestStatus `json:"status"`
	ParentId  uint              `json:"parentId"`
	ProjectId uint              `json:"projectId"`

	Ordr int `json:"ordr"`
}

func (TestCase) TableName() string {
	return "tst_cases"
}
