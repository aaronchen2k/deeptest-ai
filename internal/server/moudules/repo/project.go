package repo

import (
	"errors"
	"fmt"
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
	_logs "github.com/deeptest-com/deeptest-next/pkg/libs/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	DB        *gorm.DB `inject:""`
	*BaseRepo `inject:""`

	*UserRepo        `inject:""`
	*ProjectRoleRepo `inject:""`
}

func (r *ProjectRepo) Paginate(req v1.ReqPaginate, userId uint) (data _domain.PageData, err error) {
	var count int64
	var projectIds []uint
	r.DB.Model(&model.ProjectMember{}).
		Select("project_id").Where("user_id = ?", userId).Scan(&projectIds)

	db := r.DB.Model(&model.Project{}).Where("NOT deleted AND id IN (?)", projectIds)

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.Enabled != "" {
		db = db.Where("disabled = ?", r.IsDisable(req.Enabled))
	}

	err = db.Count(&count).Error
	if err != nil {
		_logs.Errorf("count project error", zap.String("error:", err.Error()))
		return
	}

	projects := make([]*model.Project, 0)

	err = db.
		Scopes(r.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&projects).Error
	if err != nil {
		_logs.Errorf("query project error", zap.String("error:", err.Error()))
		return
	}

	for key, project := range projects {
		user, _ := r.UserRepo.Get(project.AdminId)
		projects[key].AdminName = user.Name
	}

	data.Populate(projects, count, req.Page, req.PageSize)

	return
}

func (r *ProjectRepo) Get(id uint) (project model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("id = ?", id).
		First(&project).Error

	return
}

func (r *ProjectRepo) GetByName(projectName string, id uint) (project model.Project, err error) {
	db := r.DB.Model(&model.Project{}).
		Where("name = ? AND NOT deleted AND NOT disabled", projectName)

	if id > 0 {
		db.Where("id != ?", id)
	}

	err = db.First(&project).Error

	return
}

func (r *ProjectRepo) GetByCode(shortName string, id uint) (ret model.Project, err error) {
	db := r.DB.Model(&ret).
		Where("short_name = ? AND NOT deleted", shortName)

	if id > 0 {
		db.Where("id != ?", id)
	}
	err = db.First(&ret).Error

	return
}

func (r *ProjectRepo) GetBySpec(spec string) (project model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("spec = ?", spec).
		First(&project).Error

	return
}

func (r *ProjectRepo) Save(po *model.Project) (err error) {
	err = r.DB.Save(po).Error

	return
}

func (r *ProjectRepo) Create(req v1.ProjectReq, userId uint) (id uint, bizErr _domain.BizErr) {
	po, err := r.GetByName(req.Name, 0)
	if po.Name != "" {
		bizErr = _domain.ErrNameExist
		return
	}

	po, err = r.GetByCode(req.ShortName, 0)
	if po.ShortName != "" {
		bizErr = _domain.ErrShortNameExist
		return
	}

	// create project
	project := model.Project{ProjectBase: req.ProjectBase}
	err = r.DB.Model(&model.Project{}).Create(&project).Error
	if err != nil {
		_logs.Errorf("add project error", zap.String("error:", err.Error()))
		bizErr = _domain.SystemErr

		return
	}
	if req.AdminId != userId {
		err = r.AddProjectMember(project.ID, req.AdminId, r.BaseRepo.GetAdminRoleName())
		if err != nil {
			_logs.Errorf("添加项目角色错误", zap.String("错误:", err.Error()))
			bizErr = _domain.SystemErr
			return 0, bizErr
		}
	}
	err = r.CreateProjectRes(project.ID, userId, req.IncludeExample)

	id = project.ID

	return
}

func (r *ProjectRepo) CreateProjectRes(projectId, userId uint, IncludeExample bool) (err error) {

	// create project member
	err = r.AddProjectMember(projectId, userId, r.BaseRepo.GetAdminRoleName())
	if err != nil {
		return
	}

	return
}

func (r *ProjectRepo) Update(req v1.ProjectReq) error {
	po, _ := r.GetByName(req.Name, req.ID)
	if po.Name != "" {
		return errors.New("同名记录已存在")
	}

	po, _ = r.GetByCode(req.ShortName, req.ID)
	if po.ShortName != "" {
		return errors.New("英文缩写已存在")
	}

	project := model.Project{ProjectBase: req.ProjectBase}
	err := r.DB.Model(&model.Project{}).Where("id = ?", req.ID).Updates(&project).Error
	if err != nil {
		_logs.Errorf("update project error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *ProjectRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.Project{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		_logs.Errorf("delete project by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ProjectRepo) ListProjectByUser(userId uint) (res []model.ProjectMemberRole, err error) {
	projectRoleMap, err := r.GetProjectRoleMapByUser(userId)

	if err != nil {
		return
	}

	projectIds := make([]uint, 0)
	for k, _ := range projectRoleMap {
		projectIds = append(projectIds, k)
	}

	projects, err := r.GetProjectsByIds(projectIds)
	if err != nil {
		return
	}

	res, err = r.CombineRoleForProject(projects, projectRoleMap)

	if err != nil {
		return
	}

	//db := r.DB.Model(&model.ProjectMember{}).
	//	Joins("LEFT JOIN biz_project p ON biz_project_member.project_id=p.id").
	//	Joins("LEFT JOIN biz_project_role r ON biz_project_member.project_role_id=r.id").
	//	Select("p.*, r.id role_id, r.name role_name").
	//	Where("NOT biz_project_member.deleted")
	//
	//if !isAdminUser {
	//	db.Where("biz_project_member.user_id = ?", userId)
	//}
	//err = db.Group("biz_project_member.project_id").Find(&projects).Error
	return
}

func (r *ProjectRepo) GetProjectRoleMapByUser(userId uint) (res map[uint]uint, err error) {
	isAdminUser, err := r.UserRepo.IsAdminUser(userId)
	if err != nil {
		return
	}

	var projectMembers []model.ProjectMember
	db := r.DB.Model(&model.ProjectMember{}).
		Select("project_id, project_role_id")
	if !isAdminUser {
		db.Where("user_id = ?", userId)
	}
	if err = db.Find(&projectMembers).Error; err != nil {
		return
	}

	res = make(map[uint]uint)
	for _, v := range projectMembers {
		res[v.ProjectId] = v.ProjectRoleId
	}

	return
}

func (r *ProjectRepo) GetProjectsByIds(ids []uint) (projects []model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("id IN (?) AND NOT deleted AND NOT disabled ", ids).
		Find(&projects).Error
	return
}

func (r *ProjectRepo) CombineRoleForProject(projects []model.Project, projectRoleMap map[uint]uint) (res []model.ProjectMemberRole, err error) {
	roleIds := make([]uint, 0)
	for _, v := range projectRoleMap {
		roleIds = append(roleIds, v)
	}
	roleIds = r.ArrayRemoveUintDuplication(roleIds)

	roleIdNameMap, err := r.ProjectRoleRepo.GetRoleIdNameMap(roleIds)
	if err != nil {
		return
	}

	for _, v := range projects {
		projectMemberRole := model.ProjectMemberRole{
			Project: v,
		}
		if roleId, ok := projectRoleMap[v.ID]; ok {
			projectMemberRole.RoleId = roleId
		}
		if projectMemberRole.RoleId == 0 {
			continue
		}
		if roleName, ok := roleIdNameMap[projectMemberRole.RoleId]; ok {
			projectMemberRole.RoleName = roleName
		}
		res = append(res, projectMemberRole)
	}

	return
}

func (r *ProjectRepo) GetCurrProjectByUser(userId uint) (currProject model.Project, err error) {
	var user model.SysUser
	err = r.DB.Preload("Profile").
		Where("id = ?", userId).
		First(&user).
		Error

	err = r.DB.Model(&model.Project{}).
		Where("id = ?", user.Profile.CurrProjectId).
		First(&currProject).Error

	return
}

func (r *ProjectRepo) ListProjectsRecentlyVisited(userId uint) (projects []model.Project, err error) {
	err = r.DB.Raw(fmt.Sprintf("SELECT p.*,max( v.created_at ) visited_time FROM biz_project_recently_visited v,biz_project p WHERE v.project_id = p.id AND v.user_id = %d AND NOT p.deleted GROUP BY v.project_id ORDER BY visited_time DESC LIMIT 3", userId)).Find(&projects).Error
	return
}

func (r *ProjectRepo) ChangeProject(projectId, userId uint) (err error) {
	err = r.DB.Model(&model.SysUserProfile{}).Where("user_id = ?", userId).
		Updates(map[string]interface{}{"curr_project_id": projectId}).Error

	return
}

func (r *ProjectRepo) AddProjectMember(projectId, userId uint, role consts.RoleType) (err error) {
	var projectRole model.ProjectRole
	projectRole, err = r.ProjectRoleRepo.FindByName(role)
	if err != nil {
		return
	}

	projectMember := model.ProjectMember{UserId: userId, ProjectId: projectId, ProjectRoleId: projectRole.ID}
	err = r.DB.Create(&projectMember).Error

	return
}

func (r *ProjectRepo) Members(req v1.ReqPaginate, projectId int) (data _domain.PageData, err error) {
	req.Order = "sys_user.created_at"
	db := r.DB.Model(&model.SysUser{}).
		Select("sys_user.id, sys_user.username, sys_user.email,sys_user.name, m.project_role_id, r.name as role_name").
		Joins("left join biz_project_member m on sys_user.id=m.user_id").
		Joins("left join biz_project_role r on m.project_role_id=r.id").
		Where("m.project_id = ?", projectId)
	if req.Keywords != "" {
		db = db.Where("sys_user.username LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}

	var count int64
	err = db.Count(&count).Error
	if err != nil {
		_logs.Errorf("count users error", zap.String("error:", err.Error()))
		return
	}

	users := make([]v1.ProjectMemberResp, 0)
	err = db.
		Scopes(r.PaginateScope(req.Page, req.PageSize, "", req.Order)).
		Scan(&users).Error
	if err != nil {
		_logs.Errorf("query users error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(users, count, req.Page, req.PageSize)

	return
}

func (r *ProjectRepo) RemoveMember(userId, projectId int) (err error) {
	/*
		err = r.DB.Model(&modelRef.ProjectMember{}).
			Where("user_id = ? AND project_id = ?", userId, projectId).
			Updates(map[string]interface{}{"deleted": true}).Error
		if err != nil {
			return
		}
	*/
	err = r.DB.
		Where("user_id = ? AND project_id=?", userId, projectId).
		Delete(&model.ProjectMember{}).Error

	return
}

func (r *ProjectRepo) FindRolesByUser(userId uint) (members []model.ProjectMember, err error) {

	err = r.DB.Model(&model.ProjectMember{}).
		Joins("LEFT JOIN biz_project_role r ON biz_project_member.project_role_id=r.id").
		Select("biz_project_member.*, r.name project_role_name").
		Where("biz_project_member.user_id = ?", userId).
		Find(&members).Error

	return
}

func (r *ProjectRepo) GetProjectsAndRolesByUser(userId uint) (projectIds, roleIds []uint) {
	var members []model.ProjectMember
	r.DB.Model(&model.ProjectMember{}).
		Where("user_id = ?", userId).
		Find(&members)

	roleIdsMap := make(map[uint]uint)

	for _, member := range members {
		projectIds = append(projectIds, member.ProjectId)
		roleIdsMap[member.ProjectRoleId] = member.ProjectRoleId
	}
	for _, v := range roleIdsMap {
		roleIds = append(roleIds, v)
	}

	return
}

func (r *ProjectRepo) FindRolesByProjectAndUser(projectId, userId uint) (projectMember model.ProjectMember, err error) {
	err = r.DB.Model(&model.ProjectMember{}).
		Where("project_id = ?", projectId).
		Where("user_id = ?", userId).
		Scan(&projectMember).Error
	return
}

func (r *ProjectRepo) UpdateUserRole(req v1.UpdateProjectMemberReq) (err error) {
	err = r.DB.Model(&model.ProjectMember{}).
		Where("project_id = ?", req.ProjectId).
		Where("user_id = ?", req.UserId).
		Updates(map[string]interface{}{"project_role_id": req.ProjectRoleId}).Error

	if err != nil {
		_logs.Errorf("update project user role error", err.Error())
		return err
	}

	return
}

func (r *ProjectRepo) GetCurrProjectMemberRoleByUser(userId uint) (ret model.ProjectMember, err error) {
	curProject, err := r.GetCurrProjectByUser(userId)
	if err != nil {
		return
	}
	if curProject.ID == 0 {
		return ret, errors.New("current project is not existed")
	}
	return r.FindRolesByProjectAndUser(curProject.ID, userId)
}

func (r *ProjectRepo) GetMembersByProject(projectId uint) (ret []model.ProjectMember, err error) {
	err = r.DB.Model(&model.ProjectMember{}).
		Where("project_id = ?", projectId).
		Find(&ret).Error
	return
}

func (r *ProjectRepo) GetProjectIdsByUserIdAndRole(userId uint, roleName consts.RoleType) (projectIds []uint) {
	var projects []model.ProjectMember
	err := r.DB.Model(model.ProjectMember{}).
		Joins("LEFT JOIN biz_project_role r ON biz_project_member.project_role_id=r.id").
		Where("biz_project_member.user_id=? and r.name=? and not biz_project_member.deleted and not biz_project_member.disabled", userId, roleName).
		Find(&projects).Error
	if err != nil {
		return
	}
	for _, project := range projects {
		projectIds = append(projectIds, project.ProjectId)
	}
	return
}

func (r *ProjectRepo) ListAll() (res []model.Project, err error) {
	err = r.DB.Model(model.Project{}).
		Where("not disabled and not deleted").
		Find(&res).Error
	return
}

func (r *ProjectRepo) GetByShortName(shortName string) (project model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("short_name = ? and not deleted", shortName).
		First(&project).Error

	return
}

func (r *ProjectRepo) GetUserIdsByProjectAnRole(projectId, roleId uint) (projectMembers []model.ProjectMember, err error) {
	err = r.DB.Model(&model.ProjectMember{}).
		Where("project_id = ?", projectId).
		Where("project_role_id = ?", roleId).
		Find(&projectMembers).Error
	return
}

func (r *ProjectRepo) GetUsernamesByProjectAndRole(projectId, roleId uint, exceptUserName string) (imAccounts []string, err error) {
	conn := r.DB.Model(&model.ProjectMember{}).
		Joins("left join sys_user u on biz_project_member.user_id=u.id").
		Select("u.username")
	if projectId != 0 {
		conn = conn.Where("biz_project_member.project_id = ?", projectId)
	}
	if roleId != 0 {
		conn = conn.Where("biz_project_member.project_role_id = ?", roleId)
	}
	if exceptUserName != "" {
		conn = conn.Where("u.username != ?", exceptUserName)
	}
	err = conn.Where("not biz_project_member.deleted and not u.deleted").Find(&imAccounts).Error
	return
}

func (r *ProjectRepo) ListByUsername(username string) (res []model.Project, err error) {
	err = r.DB.Model(model.Project{}).
		Joins("LEFT JOIN biz_project_member m ON biz_project.id=m.project_id").
		Joins("LEFT JOIN sys_user u ON m.user_id=u.id").
		Where("u.username = ?", username).
		Where("not biz_project.disabled and not biz_project.deleted and not m.disabled and not m.deleted").
		Find(&res).Error
	return
}

func (r *ProjectRepo) BatchGetByShortNames(shortNames []string) (ret []model.Project, err error) {
	err = r.DB.Model(&ret).
		Where("short_name IN (?) AND NOT deleted", shortNames).
		Find(&ret).Error

	return
}

func (r *ProjectRepo) AddMemberIfNotExisted(projectId, userId uint, role consts.RoleType) (err error) {
	isMember, err := r.IfProjectMember(userId, projectId)
	if err != nil || isMember {
		return
	}

	err = r.AddProjectMember(projectId, userId, role)
	return
}
func (r *ProjectRepo) IfProjectMember(userId, projectId uint) (res bool, err error) {
	var count int64
	err = r.DB.Model(&model.ProjectMember{}).Where("user_id=? and project_id=?", userId, projectId).Count(&count).Error
	if err != nil {
		return
	}
	res = count > 0
	return
}

func (r *ProjectRepo) FindRolesByProjectsAndUsername(username string, projectIds []uint) (members []model.ProjectMember, err error) {
	err = r.DB.Model(&model.ProjectMember{}).
		Joins("LEFT JOIN biz_project_role r ON biz_project_member.project_role_id=r.id").
		Joins("LEFT JOIN sys_user u ON biz_project_member.user_id=u.id").
		Select("biz_project_member.*, r.name project_role_name").
		Where("u.username = ?", username).
		Where("biz_project_member.project_id IN (?)", projectIds).
		Find(&members).Error

	return
}

func (r *ProjectRepo) GetUserProjectRoleMap(username string, projectIds []uint) (res map[uint]consts.RoleType, err error) {
	projectRoles, err := r.FindRolesByProjectsAndUsername(username, projectIds)
	if err != nil {
		return
	}

	res = make(map[uint]consts.RoleType)
	for _, v := range projectRoles {
		res[v.ProjectId] = v.ProjectRoleName
	}

	return
}

func (r *ProjectRepo) GetProjectMemberCount(projectId uint) (count int64, err error) {
	err = r.DB.Model(&model.ProjectMember{}).Where("project_id=?", projectId).Count(&count).Error
	if err != nil {
		return
	}
	return
}

func (r *ProjectRepo) GetProjectMemberList(projectId uint) (list []model.ProjectMember, err error) {
	err = r.DB.Model(&model.ProjectMember{}).Where("project_id=?", projectId).Find(&list).Error
	if err != nil {
		return
	}
	return
}
