package service

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/jinzhu/copier"
)

type UserService struct {
	UserRepo *repo.UserRepo `inject:""`
}

func (s *UserService) Paginate(req v1.UserPageReq) (ret _domain.PageData, err error) {
	ret, err = s.UserRepo.Paginate(req)

	return
}

func (s *UserService) Get(id int) (ret v1.UserResp, err error) {
	po, err := s.UserRepo.Get(uint(id))
	if err != nil {
		return
	}

	copier.CopyWithOption(&ret, po, copier.Option{DeepCopy: true})

	return
}

func (s *UserService) Create(req *v1.UserReq) (po model.SysUser, err error) {
	po, err = s.UserRepo.Create(req)
	return
}

func (s *UserService) Update(req *v1.UserReq) (err error) {
	po := model.SysUser{}

	copier.CopyWithOption(&po, req, copier.Option{DeepCopy: true})

	err = s.UserRepo.Update(po)
	return
}

func (s *UserService) Delete(id int) (err error) {
	err = s.UserRepo.Delete(uint(id))
	return
}
