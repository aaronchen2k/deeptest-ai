package handler

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type BaseCtrl struct {
}

func (c *BaseCtrl) getTenantId(ctx iris.Context) domain.TenantId {
	return GetTenantId(ctx)
}

func GetTenantId(ctx *context.Context) (ret domain.TenantId) {
	tenantId := ctx.GetHeader("tenantId")

	ret = domain.TenantId(tenantId)

	return
}
