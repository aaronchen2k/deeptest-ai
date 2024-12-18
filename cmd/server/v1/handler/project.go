package handler

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/auth/iris"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/service"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/kataras/iris/v12"
)

type ProjectCtrl struct {
	ProjectService *service.ProjectService `inject:""`
	BaseCtrl
}

func (c *ProjectCtrl) List(ctx iris.Context) {
	userId := multi_iris.GetUserId(ctx)

	var req v1.ReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}
	req.ConvertParams()

	data, err := c.ProjectService.Paginate(req, userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: data})
}

func (c *ProjectCtrl) Get(ctx iris.Context) {
	var req _domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}
	project, err := c.ProjectService.Get(req.Id)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: project})
}

func (c *ProjectCtrl) Create(ctx iris.Context) {
	userId := multi_iris.GetUserId(ctx)

	req := v1.ProjectReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	id, bizErr := c.ProjectService.Create(req, userId)
	if bizErr.Code != 0 {
		ctx.JSON(_domain.Response{Code: bizErr.Code, Data: nil, Msg: bizErr.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: iris.Map{"id": id}})
}

func (c *ProjectCtrl) Update(ctx iris.Context) {
	var req v1.ProjectReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	err = c.ProjectService.Update(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}

func (c *ProjectCtrl) Delete(ctx iris.Context) {
	var req _domain.ReqId
	err := ctx.ReadParams(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	err = c.ProjectService.DeleteById(req.Id)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}

func (c *ProjectCtrl) GetByUser(ctx iris.Context) {
	userId := multi_iris.GetUserId(ctx)

	projects, currProject, recentProjects, err := c.ProjectService.GetByUser(userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ret := iris.Map{"projects": projects, "currProject": currProject, "recentProjects": recentProjects}
	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: ret})
}

func (c *ProjectCtrl) ChangeProject(ctx iris.Context) {
	userId := multi_iris.GetUserId(ctx)

	req := v1.ProjectReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	err = c.ProjectService.ChangeProject(req.ID, userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	projects, currProject, recentProjects, err := c.ProjectService.GetByUser(userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ret := iris.Map{"projects": projects, "currProject": currProject, "recentProjects": recentProjects}
	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: ret})
}

func (c *ProjectCtrl) Members(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("currProjectId")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	var req v1.ReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}
	req.ConvertParams()

	data, err := c.ProjectService.Members(req, projectId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: data})
}

func (c *ProjectCtrl) RemoveMember(ctx iris.Context) {
	req := v1.ProjectMemberRemoveReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	err = c.ProjectService.RemoveMember(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}

func (c *ProjectCtrl) ChangeUserRole(ctx iris.Context) {
	req := v1.UpdateProjectMemberReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	err = c.ProjectService.UpdateMemberRole(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}
