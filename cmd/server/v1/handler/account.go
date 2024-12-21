package handler

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/auth"
	multi_iris "github.com/deeptest-com/deeptest-next/internal/pkg/core/auth/iris"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/service"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/kataras/iris/v12"
)

type AccountCtrl struct {
	BaseCtrl
	AccountService *service.AccountService `inject:""`
	UserService    *service.UserService    `inject:""`
}

func (c *AccountCtrl) Login(ctx iris.Context) {
	req := &v1.LoginReq{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	token, id, err := c.AccountService.GetAccessToken(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: iris.Map{"accessToken": token, "user": iris.Map{"id": id}}})
}

func (c AccountCtrl) Logout(ctx iris.Context) {
	token := multi_iris.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(_domain.Response{Code: _domain.Success.Code})
		return
	}

	err := c.AccountService.DeleteToken(string(token))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}

func (c AccountCtrl) CleanToken(ctx iris.Context) {
	token := multi_iris.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: "授权凭证为空"})
		return
	}

	if err := c.AccountService.CleanToken(auth.AdminAuthority, string(token)); err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}

func (c AccountCtrl) Info(ctx iris.Context) {
	userId := multi_iris.GetUserId(ctx)

	user, err := c.AccountService.GetInfo(userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: user})
}

func (c AccountCtrl) Codes(ctx iris.Context) {
	userId := multi_iris.GetUserId(ctx)

	codes, err := c.AccountService.GetCodes(userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: codes})
}

func (c *AccountCtrl) UpdateUserProject(ctx iris.Context) {
	userId := multi_iris.GetUserId(ctx)

	req := &v1.UpdateUserProjectReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.AccountService.UpdateUserProject(req, userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}
