package handler

import (
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/kataras/iris/v12"
)

type BaseCtrl struct {
}

func (c *BaseCtrl) successResp() (ret _domain.Response) {
	ret = _domain.Response{Code: _domain.Success.Code, Data: iris.Map{"success": true}}

	return
}

//func (c *BaseCtrl) getTenantId(ctx iris.Context) domain.TenantId {
//	return GetTenantId(ctx)
//}
//
//func GetTenantId(ctx *context.Context) (ret domain.TenantId) {
//	tenantId := ctx.GetHeader("tenantId")
//
//	ret = domain.TenantId(tenantId)
//
//	return
//}
