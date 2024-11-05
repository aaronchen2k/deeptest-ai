package v1

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
)

type RolePageReq struct {
	_domain.PageReq
	Name string `json:"name"`
}
type RolePageResp struct {
	Item []*RoleResp
}

type RoleReq struct {
	BaseDomain
	domain.BaseRole

	Perms [][]string `json:"perms"`
}
type RoleResp struct {
	BaseDomain
	domain.BaseRole
}
