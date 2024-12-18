package handler

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/service"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/str"
)

type DataCtrl struct {
	BaseCtrl
	DataService *service.DataService `inject:""`
}

func (c *DataCtrl) Init(ctx iris.Context) {
	req := v1.DataReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	err := c.DataService.InitDB(req)

	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}

func (c *DataCtrl) Check(ctx iris.Context) {
	if c.DataService.DataRepo.DB == nil {
		ctx.JSON(_domain.Response{Code: _domain.NeedInitErr.Code, Data: iris.Map{
			"needInit": true,
		}, Msg: str.Join(_domain.NeedInitErr.Msg, ":数据库初始化失败")})
		return

	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: iris.Map{
		"needInit": false,
	}, Msg: _domain.Success.Msg})
}
