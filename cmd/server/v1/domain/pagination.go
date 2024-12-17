package v1

import (
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
)

type ReqPaginate struct {
	_domain.PaginateReq

	Keywords string `json:"keywords"`
	Status   string `json:"status"`
	Enabled  string `json:"enabled"`

	ProjectId uint `json:"projectId"`
}
