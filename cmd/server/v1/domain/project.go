package v1

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
)

type ProjectReq struct {
	BaseDomain
	domain.ProjectBase
}

type ProjectResp struct {
	_domain.PaginateReq
	domain.ProjectBase
}

type ProjectMemberResp struct {
	Id            uint            `json:"id"`
	Username      string          `json:"username"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	RoleName      consts.RoleType `json:"roleName"`
	ProjectRoleId uint            `json:"roleId"`
}

type UpdateProjectMemberReq struct {
	ProjectId     uint `json:"projectId"`
	ProjectRoleId uint `json:"projectRoleId"`
	UserId        uint `json:"userId"`
}

type ProjectMemberRemoveReq struct {
	UserId    int `json:"userId"`
	ProjectId int `json:"projectId"`
}
