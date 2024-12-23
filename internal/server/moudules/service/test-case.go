package service

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
)

type CaseService struct {
	CaseRepo *repo.CaseRepo `inject:""`
}

func (s *CaseService) LoadTree(projectId int) (ret []*domain.CaseNode, err error) {
	root, err := s.CaseRepo.LoadTree(uint(projectId))

	if root != nil {
		ret = root.Children
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
