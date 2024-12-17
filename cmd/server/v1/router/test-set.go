package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type SetModule struct {
	SetCtrl *handler.SetCtrl `inject:""`
}

// Party 集合
func (m *SetModule) Party() func(index iris.Party) {
	return func(party iris.Party) {
		party.Use(middleware.MultiHandler(), middleware.Casbin())

		party.Get("/", m.SetCtrl.List).Name = "集合列表"
		party.Get("/{id:uint}", m.SetCtrl.Get).Name = "集合详情"
		party.Post("/", m.SetCtrl.Create).Name = "新建集合"
		party.Put("/", m.SetCtrl.Update).Name = "更新集合"
		party.Delete("/{id:uint}", m.SetCtrl.Delete).Name = "删除集合"
	}
}
