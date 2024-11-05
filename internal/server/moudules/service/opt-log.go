package service

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
)

type OplogService struct {
	OplogRepo *repo.OplogRepo `inject:""`
}

func (s *OplogService) Paginate(req v1.OptLogPageReq) (ret _domain.PageData, err error) {
	ret, err = s.OplogRepo.Paginate(req)

	return
}

func (s *OplogService) List() (ret []model.SysOplog, err error) {
	ret, err = s.OplogRepo.List()

	return
}
