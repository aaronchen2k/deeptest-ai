package service

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
)

type SetService struct {
	SetRepo *repo.SetRepo `inject:""`
}

func (s *SetService) Paginate(req v1.ReqPaginate, projectId int) (ret _domain.PageData, err error) {
	ret, err = s.SetRepo.Paginate(req, projectId)

	if err != nil {
		return
	}

	return
}

func (s *SetService) GetById(id uint) (ret model.TestSet, err error) {
	ret, err = s.SetRepo.Get(id)
	if err != nil {
		return
	}

	return
}

func (s *SetService) Create(req model.TestSet) (po model.TestSet, err error) {
	po, err = s.SetRepo.Create(req)

	return
}

func (s *SetService) Update(req model.TestSet) (err error) {
	err = s.SetRepo.Update(req)

	return
}

func (s *SetService) DeleteById(id uint) (err error) {
	err = s.SetRepo.Delete(id)

	return
}
