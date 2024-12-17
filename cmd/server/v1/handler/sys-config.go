package handler

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/service"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/kataras/iris/v12"
	"os"
)

type ConfigCtrl struct {
	BaseCtrl
	ConfigService *service.SettingsService `inject:""`
}

const token = "a1bc**2d&&423qvdw"

func (c *ConfigCtrl) Get(ctx iris.Context) {
	data := iris.Map{
		"demoTestSite": os.Getenv("DEMO_TEST_SITE"),
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: data})
}

func (c *ConfigCtrl) GetValue(ctx iris.Context) {
	key := ctx.URLParam("key")
	if key == "" {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	value, err := c.ConfigService.Get(key)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: value})
}

func (c *ConfigCtrl) Save(ctx iris.Context) {
	headerToken := ctx.Request().Header.Get("token")
	if headerToken != token {
		ctx.JSON(_domain.Response{Code: _domain.AuthActionErr.Code, Msg: _domain.AuthActionErr.Msg})
		return
	}
	req := model.Settings{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	err = c.ConfigService.Save(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}
