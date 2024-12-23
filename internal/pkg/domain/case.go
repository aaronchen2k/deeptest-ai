package domain

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/kataras/iris/v12"
)

type CaseNode struct {
	Id int64 `json:"id"`

	Title string              `json:"title"`
	Desc  string              `json:"desc,omitempty"`
	Type  consts.TreeNodeType `json:"type"`
	IsDir bool                `json:"isDir"`

	ParentId  int64 `json:"parentId"`
	ProjectId uint  `json:"projectId"`
	UseID     uint  `json:"useId"`

	Ordr     int         `json:"ordr"`
	Children []*CaseNode `json:"children,omitempty"`
	Slots    iris.Map    `json:"slots"`
	Count    int         `json:"count,omitempty"`
}
