package _domain

import (
	_consts "github.com/deeptest-com/deeptest-next/pkg/libs/consts"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
)

type Model struct {
	Id        uint   `json:"id"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type ReqId struct {
	Id uint `json:"id" param:"id"`
}

type PaginateReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Field    string `json:"field"`
	Order    string `json:"order"`
}

func (r *PaginateReq) ConvertParams() {
	r.Field = _str.SnakeCase(r.Field)
	r.Order = _consts.SortMap[r.Order]
}
