package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	middleware2 "github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type OptLogModule struct {
	OptLogCtrl *handler.OplogCtrl `inject:""`
}

func (m *OptLogModule) Party() func(index iris.Party) {
	return func(index iris.Party) {
		index.Use(middleware2.MultiHandler(), middleware2.Casbin())

		index.Get("/", m.OptLogCtrl.Paginate).Name = "操作日志列表"
	}
}
