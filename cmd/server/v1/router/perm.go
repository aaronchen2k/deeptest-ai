package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	middleware2 "github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type PermModule struct {
	PermCtrl *handler.PermCtrl `inject:""`
}

func (m *PermModule) Party() func(index iris.Party) {
	return func(index iris.Party) {
		index.Use(middleware2.JwtHandler(), middleware2.Casbin())

		index.Get("/", m.PermCtrl.Paginate).Name = "权限列表"
		index.Get("/{id:uint}", m.PermCtrl.Get).Name = "权限详情"
	}
}
