package v1

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
)

type PermPageReq struct {
	_domain.PageReq
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}
type PermPageResp struct {
	Item []*PermResp
}

type PermReq struct {
	BaseDomain
	domain.BasePermission
}
type PermResp struct {
	BaseDomain
	domain.BasePermission
}
