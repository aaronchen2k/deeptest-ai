package service

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/jinzhu/copier"
)

type PermService struct {
	PermRepo *repo.PermRepo `inject:""`
}

func (s *PermService) Paginate(req v1.PermPageReq) (ret _domain.PageData, err error) {
	ret, err = s.PermRepo.Paginate(req)

	return
}

func (s *PermService) Get(id int) (ret v1.PermPageResp, err error) {
	po, err := s.PermRepo.Get(uint(id))
	if err != nil {
		return
	}

	copier.CopyWithOption(&ret, po, copier.Option{DeepCopy: true})

	return
}
