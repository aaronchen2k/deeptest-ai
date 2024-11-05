package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	middleware2 "github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type AccountModule struct {
	AccountCtrl *handler.AccountCtrl `inject:""`
}

func (m *AccountModule) Party() func(public iris.Party) {
	return func(public iris.Party) {
		public.Post("/login", m.AccountCtrl.Login).Name = "登录"

		public.Use(middleware2.MultiHandler(), middleware2.Casbin())
		public.Get("/logout", m.AccountCtrl.Logout).Name = "退出"
		public.Get("/clear", m.AccountCtrl.CleanToken).Name = "清空Token"
	}
}
