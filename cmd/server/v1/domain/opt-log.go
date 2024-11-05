package v1

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
)

type OptLogPageReq struct {
	_domain.PageReq
	Name     string `json:"name"`
	UserName string `json:"userName"`
}

type OptLogPageResp struct {
	Item []*OptLogResp
}

type OptLogResp struct {
	domain.BaseOplog
}
