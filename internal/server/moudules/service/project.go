package service

import (
	"fmt"
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
)

type ProjectService struct {
	ProjectRepo     *repo.ProjectRepo     `inject:""`
	ProjectRoleRepo *repo.ProjectRoleRepo `inject:""`
	UserRepo        *repo.UserRepo        `inject:""`
}

func (s *ProjectService) Load(userId uint) (curr v1.ProjectReq, items []v1.ProjectReq, err error) {
	curr, items, err = s.ProjectRepo.Load(userId)
	if err != nil {
		return
	}

	return
}

func (s *ProjectService) Paginate(req v1.ReqPaginate, userId uint) (ret _domain.PageData, err error) {
	ret, err = s.ProjectRepo.Paginate(req, userId)
	if err != nil {
		return
	}

	return
}

func (s *ProjectService) Get(id uint) (model.Project, error) {
	return s.ProjectRepo.Get(id)
}

func (s *ProjectService) Create(req v1.ProjectReq, userId uint) (id uint, err _domain.BizErr) {
	id, err = s.ProjectRepo.Create(req, userId)
	if err.Code != 0 {
		return
	}

	return
}

func (s *ProjectService) Update(req v1.ProjectReq, userId uint) (err error) {
	err = s.ProjectRepo.Update(req, userId)
	if err != nil {
		return
	}

	return
}

func (s *ProjectService) DeleteById(id uint) error {
	return s.ProjectRepo.DeleteById(id)
}

func (s *ProjectService) GetByUser(userId uint) (projects []model.ProjectMemberRole, currProject model.Project, recentProjects []model.Project, err error) {
	projects, err = s.ProjectRepo.ListProjectByUser(userId)
	currProject, err = s.ProjectRepo.GetCurrProjectByUser(userId)
	recentProjects, err = s.ProjectRepo.ListProjectsRecentlyVisited(userId)

	return
}

func (s *ProjectService) ChangeProject(projectId, userId uint) (err error) {
	err = s.ProjectRepo.ChangeProject(projectId, userId)

	return
}

func (s *ProjectService) Members(req v1.ReqPaginate, projectId int) (data _domain.PageData, err error) {
	data, err = s.ProjectRepo.Members(req, projectId)

	return
}

func (s *ProjectService) RemoveMember(req v1.ProjectMemberRemoveReq) (err error) {
	err = s.ProjectRepo.RemoveMember(req.UserId, req.ProjectId)

	return
}

func (s *ProjectService) UpdateMemberRole(req v1.UpdateProjectMemberReq) (err error) {
	return s.ProjectRepo.UpdateUserRole(req)
}

func (s *ProjectService) GetCurrProjectByUser(userId uint) (currProject model.Project, err error) {
	currProject, err = s.ProjectRepo.GetCurrProjectByUser(userId)

	return
}

func (s *ProjectService) AllProjectList(username string) (res []model.Project, err error) {
	return s.ProjectRepo.ListByUsername(username)
}

func (s *ProjectService) GetProjectRole(username, projectCode string) (role string, err error) {
	var user model.SysUser
	user, _ = s.UserRepo.GetByUserName(username)
	if user.ID == 0 {
		err = fmt.Errorf("用户名不存在")
		return
	}
	var project model.Project
	project, _ = s.ProjectRepo.GetByShortName(projectCode)
	if project.ID == 0 {
		err = fmt.Errorf("项目不存在")
		return
	}

	var projectRole model.ProjectRole
	projectRole, _ = s.ProjectRoleRepo.ProjectUserRoleList(user.ID, project.ID)
	if projectRole.Name == "" {
		err = fmt.Errorf("用户角色不存在")
		return
	}

	return string(projectRole.Name), err
}
