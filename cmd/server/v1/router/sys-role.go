package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type RoleModule struct {
	RoleCtrl *handler.RoleCtrl `inject:""`
}

func (m *RoleModule) Party() func(index iris.Party) {
	return func(index iris.Party) {
		index.Use(middleware.JwtHandler(), middleware.Casbin())

		index.Get("/", m.RoleCtrl.Paginate).Name = "角色列表"
		index.Get("/{id:uint}", m.RoleCtrl.Get).Name = "角色详情"
	}
}
