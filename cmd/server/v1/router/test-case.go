package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type CaseModule struct {
	CaseCtrl *handler.CaseCtrl `inject:""`
}

// Party 项目
func (m *CaseModule) Party() func(index iris.Party) {
	return func(party iris.Party) {
		party.Use(middleware.JwtHandler(), middleware.Casbin())

		party.Get("/", m.CaseCtrl.List).Name = "用例列表"
		party.Get("/{id:uint}", m.CaseCtrl.Get).Name = "用例详情"
		party.Post("/", m.CaseCtrl.Create).Name = "新建用例"
		party.Put("/", m.CaseCtrl.Update).Name = "更新用例"
		party.Delete("/{id:uint}", m.CaseCtrl.Delete).Name = "删除用例"
	}
}
