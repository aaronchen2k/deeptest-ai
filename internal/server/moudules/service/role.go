package service

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/jinzhu/copier"
)

type RoleService struct {
	RoleRepo *repo.RoleRepo `inject:""`
}

func (s *RoleService) Paginate(req v1.RolePageReq) (ret _domain.PageData, err error) {
	ret, err = s.RoleRepo.Paginate(req)

	return
}

func (s *RoleService) Get(id int) (ret v1.RoleResp, err error) {
	po, err := s.RoleRepo.Get(uint(id))
	if err != nil {
		return
	}

	copier.CopyWithOption(&ret, po, copier.Option{DeepCopy: true})

	return
}
