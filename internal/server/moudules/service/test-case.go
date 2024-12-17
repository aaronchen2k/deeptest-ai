package service

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
)

type CaseService struct {
	CaseRepo *repo.CaseRepo `inject:""`
}

func (s *CaseService) Paginate(req v1.ReqPaginate, projectId int) (ret _domain.PageData, err error) {
	ret, err = s.CaseRepo.Paginate(req, projectId)

	if err != nil {
		return
	}

	return
}

func (s *CaseService) GetById(id uint) (ret model.TestCase, err error) {
	ret, err = s.CaseRepo.Get(id)
	if err != nil {
		return
	}

	return
}

func (s *CaseService) Create(req model.TestCase) (po model.TestCase, err error) {
	po, err = s.CaseRepo.Create(req)

	return
}

func (s *CaseService) Update(req model.TestCase) (err error) {
	err = s.CaseRepo.Update(req)

	return
}

func (s *CaseService) DeleteById(id uint) (err error) {
	err = s.CaseRepo.Delete(id)

	return
}
