package v1

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
)

type UserPageReq struct {
	_domain.PageReq
	Name string `json:"name"`
}
type UserPageResp struct {
	Item []*UserResp
}

type UserReq struct {
	BaseDomain
	domain.BaseUser

	Password  string   `json:"password"`
	RoleNames []string `json:"role_ids"`
}
type UserResp struct {
	BaseDomain
	domain.BaseUser

	Roles []string `gorm:"-" json:"roles"`
}
