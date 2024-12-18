package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type PlanModule struct {
	PlanCtrl *handler.PlanCtrl `inject:""`
}

// Party 项目
func (m *PlanModule) Party() func(index iris.Party) {
	return func(party iris.Party) {
		party.Use(middleware.JwtHandler(), middleware.Casbin())

		party.Get("/", m.PlanCtrl.List).Name = "项目列表"
		party.Get("/{id:uint}", m.PlanCtrl.Get).Name = "计划详情"
		party.Post("/", m.PlanCtrl.Create).Name = "新建计划"
		party.Put("/", m.PlanCtrl.Update).Name = "更新计划"
		party.Delete("/{id:uint}", m.PlanCtrl.Delete).Name = "删除计划"
	}
}
