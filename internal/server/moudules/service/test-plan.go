package service

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
)

type PlanService struct {
	PlanRepo *repo.PlanRepo `inject:""`
}

func (s *PlanService) Paginate(req v1.ReqPaginate, projectId int) (ret _domain.PageData, err error) {
	ret, err = s.PlanRepo.Paginate(req, projectId)

	if err != nil {
		return
	}

	return
}

func (s *PlanService) GetById(id uint) (ret model.TestPlan, err error) {
	ret, err = s.PlanRepo.Get(id)
	if err != nil {
		return
	}

	return
}

func (s *PlanService) Create(req model.TestPlan) (po model.TestPlan, err error) {
	po, err = s.PlanRepo.Create(req)

	return
}

func (s *PlanService) Update(req model.TestPlan) (err error) {
	err = s.PlanRepo.Update(req)

	return
}

func (s *PlanService) DeleteById(id uint) (err error) {
	err = s.PlanRepo.Delete(id)

	return
}
